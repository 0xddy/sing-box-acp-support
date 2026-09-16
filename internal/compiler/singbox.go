package compiler

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

type SingBoxConfig struct {
	Log          map[string]any   `json:"log,omitempty"`
	DNS          map[string]any   `json:"dns,omitempty"`
	HTTPClients  []map[string]any `json:"http_clients,omitempty"`
	Inbounds     []map[string]any `json:"inbounds"`
	Outbounds    []map[string]any `json:"outbounds"`
	Route        map[string]any   `json:"route,omitempty"`
	Experimental map[string]any   `json:"experimental,omitempty"`
}

func Compile(top topology.MachineTopology) ([]byte, error) {
	if top.MachineID == "" {
		return nil, errors.New("machine_id is required")
	}

	nodes := append([]topology.NodeInstance(nil), top.Nodes...)
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].NodeID < nodes[j].NodeID
	})

	inbounds := make([]map[string]any, 0, len(nodes))
	compiledNodes := make([]compiledInboundMetadata, 0, len(nodes))
	for _, node := range nodes {
		compiled, err := compileInbound(node)
		if err != nil {
			return nil, err
		}
		inbounds = append(inbounds, compiled.Inbound)
		compiledNodes = append(compiledNodes, compiledInboundMetadata{Tag: compiled.Tag, Sniff: compiled.Sniff})
	}

	outbounds, outboundTags, defaultRouteFinal, err := compileOutbounds(top.Outbounds)
	if err != nil {
		return nil, err
	}
	route, err := compileRoute(top.Route, outboundTags, defaultRouteFinal)
	if err != nil {
		return nil, err
	}
	route = prependRouteRules(route, sniffRouteRules(compiledNodes))
	httpClients := []map[string]any(nil)
	if needsDefaultRuleSetHTTPClient(top.Route) {
		const clientTag = "__acp_rule_set_download"
		client := map[string]any{"tag": clientTag}
		defaultOutboundTag, _ := route["final"].(string)
		defaultOutbound := outbounds[0]
		if defaultOutboundTag != "" {
			for _, outbound := range outbounds {
				if outbound["tag"] == defaultOutboundTag {
					defaultOutbound = outbound
					break
				}
			}
		}
		if defaultOutbound["type"] == topology.OutboundTypeDirect {
			// An empty Direct outbound cannot be used as a detour. Its dialer
			// options can instead be used directly by the shared HTTP client.
			for key, value := range defaultOutbound {
				if key != "type" && key != "tag" && key != "proxy_protocol" {
					client[key] = value
				}
			}
		} else {
			client["detour"] = defaultOutbound["tag"]
		}
		httpClients = []map[string]any{client}
		route["default_http_client"] = clientTag
	}
	dns, err := compileDNS(top.DNS)
	if err != nil {
		return nil, err
	}

	cfg := SingBoxConfig{
		Log: map[string]any{
			"level":     "info",
			"timestamp": true,
		},
		DNS:         dns,
		HTTPClients: httpClients,
		Inbounds:    inbounds,
		Outbounds:   outbounds,
		Route:       route,
		Experimental: map[string]any{
			"cache_file": map[string]any{
				"enabled": true,
				"path":    "runtime/cache.db",
			},
		},
	}
	return json.MarshalIndent(cfg, "", "  ")
}

func needsDefaultRuleSetHTTPClient(route *topology.Route) bool {
	if route == nil {
		return false
	}
	for _, ruleSet := range route.RuleSets {
		if ruleSet.Type == "remote" && ruleSet.DownloadDetour == "" {
			return true
		}
	}
	return false
}

type compiledInboundMetadata struct {
	Tag   string
	Sniff bool
}

func sniffRouteRules(nodes []compiledInboundMetadata) []map[string]any {
	rules := make([]map[string]any, 0)
	for _, node := range nodes {
		if !node.Sniff {
			continue
		}
		tag := node.Tag
		if tag == "" {
			continue
		}
		rules = append(rules, map[string]any{
			"inbound": []string{tag},
			"action":  "sniff",
			"timeout": "300ms",
		})
	}
	return rules
}

func prependRouteRules(route map[string]any, rules []map[string]any) map[string]any {
	if len(rules) == 0 {
		return route
	}
	if route == nil {
		route = map[string]any{}
	}
	existing, _ := route["rules"].([]map[string]any)
	merged := make([]map[string]any, 0, len(rules)+len(existing))
	merged = append(merged, rules...)
	merged = append(merged, existing...)
	route["rules"] = merged
	return route
}
