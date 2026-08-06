package compiler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func compileDNS(dns *topology.DNS) (map[string]any, error) {
	compiled := map[string]any{
		"servers": defaultDNSServers(),
		"final":   "default-dns",
	}
	if dns == nil {
		return compiled, nil
	}
	if len(dns.Servers) > 0 {
		servers := make([]map[string]any, 0, len(dns.Servers))
		for _, server := range dns.Servers {
			compiledServer, err := compileDNSServer(server)
			if err != nil {
				return nil, err
			}
			servers = append(servers, compiledServer)
		}
		compiled["servers"] = servers
	}
	if dns.Final != "" {
		compiled["final"] = dns.Final
	}
	if len(dns.Rules) == 0 {
		return compiled, nil
	}
	rules := make([]map[string]any, 0, len(dns.Rules))
	for _, rule := range dns.Rules {
		compiledRule, err := compileDNSRule(rule)
		if err != nil {
			return nil, err
		}
		rules = append(rules, compiledRule)
	}
	compiled["rules"] = rules
	return compiled, nil
}

func defaultDNSServers() []map[string]any {
	return []map[string]any{
		{
			"type":   "https",
			"tag":    "default-dns",
			"server": "1.1.1.1",
		},
	}
}

func compileDNSServer(server topology.DNSServer) (map[string]any, error) {
	if server.Type == "" {
		return nil, errors.New("dns server type is required")
	}
	if server.Tag == "" {
		return nil, errors.New("dns server tag is required")
	}
	if server.Server == "" {
		return nil, errors.New("dns server address is required")
	}
	compiled := map[string]any{
		"type":   server.Type,
		"tag":    server.Tag,
		"server": server.Server,
	}
	applyOptionalString(compiled, "detour", server.Detour)
	return compiled, nil
}

func compileDNSRule(rule topology.DNSRule) (map[string]any, error) {
	if rule.Action == "" {
		return nil, errors.New("dns rule action is required")
	}
	compiled := map[string]any{"action": rule.Action}
	applyStringList(compiled, "inbound", rule.Inbound)
	applyStringList(compiled, "domain", rule.Domain)
	applyStringList(compiled, "domain_suffix", rule.DomainSuffix)
	applyStringList(compiled, "domain_keyword", rule.DomainKeyword)
	applyStringList(compiled, "domain_regex", rule.DomainRegex)
	applyStringList(compiled, "rule_set", rule.RuleSet)
	applyOptionalString(compiled, "rcode", rule.RCode)
	applyOptionalString(compiled, "server", rule.Server)
	applyOptionalString(compiled, "method", rule.Method)
	if rule.NoDrop {
		compiled["no_drop"] = rule.NoDrop
	}
	applyStringList(compiled, "answer", rule.Answer)
	applyStringList(compiled, "ns", rule.NS)
	applyStringList(compiled, "extra", rule.Extra)
	if rule.DisableCache {
		compiled["disable_cache"] = rule.DisableCache
	}
	if rewriteTTL := strings.TrimSpace(rule.RewriteTTL); rewriteTTL != "" {
		value, err := strconv.ParseUint(rewriteTTL, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("dns rule rewrite_ttl %q must be a uint32: %w", rule.RewriteTTL, err)
		}
		compiled["rewrite_ttl"] = uint32(value)
	}
	applyOptionalString(compiled, "client_subnet", rule.ClientSubnet)
	return compiled, nil
}
