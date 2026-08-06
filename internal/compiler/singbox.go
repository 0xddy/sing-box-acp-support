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
	dns, err := compileDNS(top.DNS)
	if err != nil {
		return nil, err
	}

	cfg := SingBoxConfig{
		Log: map[string]any{
			"level":     "info",
			"timestamp": true,
		},
		DNS:       dns,
		Inbounds:  inbounds,
		Outbounds: outbounds,
		Route:     route,
		Experimental: map[string]any{
			"cache_file": map[string]any{
				"enabled": true,
				"path":    "runtime/cache.db",
			},
		},
	}
	return json.MarshalIndent(cfg, "", "  ")
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
