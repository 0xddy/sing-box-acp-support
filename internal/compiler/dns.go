package compiler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	C "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
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
		servers, err := compileDNSServers(dns.Servers, dns.Final)
		if err != nil {
			return nil, err
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
	return []map[string]any{defaultDNSServer("default-dns")}
}

func defaultDNSServer(tag string) map[string]any {
	return map[string]any{
		"type":   "https",
		"tag":    tag,
		"server": "1.1.1.1",
	}
}

func compileDNSServers(servers []topology.DNSServer, final string) ([]map[string]any, error) {
	compiled := make([]map[string]any, 0, len(servers))
	tags := make(map[string]struct{}, len(servers))
	var needsBootstrap []int
	bootstrapTag := ""
	for _, server := range servers {
		entry, err := compileDNSServer(server)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, entry)
		tags[server.Tag] = struct{}{}
		// A detoured DNS server delegates its destination to the outbound.
		// Adding a local resolver there could change the selected egress or
		// introduce a cycle through that outbound's own domain resolver.
		if server.Detour != "" {
			continue
		}
		switch server.Type {
		case C.DNSTypeUDP, C.DNSTypeTCP, C.DNSTypeTLS, C.DNSTypeHTTPS, C.DNSTypeQUIC:
			if (option.DNSServerAddressOptions{Server: server.Server}).ServerIsDomain() {
				needsBootstrap = append(needsBootstrap, len(compiled)-1)
			} else if bootstrapTag == "" || server.Tag == final {
				// A literal address without a detour has no DNS dependency, so
				// it can safely bootstrap any of the hostname servers.
				bootstrapTag = server.Tag
			}
		}
	}
	if len(needsBootstrap) == 0 {
		return compiled, nil
	}
	if bootstrapTag == "" {
		const baseTag = "__acp_dns_bootstrap"
		bootstrapTag = baseTag
		for suffix := 1; ; suffix++ {
			if _, exists := tags[bootstrapTag]; !exists {
				break
			}
			bootstrapTag = fmt.Sprintf("%s_%d", baseTag, suffix)
		}
		compiled = append(compiled, defaultDNSServer(bootstrapTag))
	}
	for _, index := range needsBootstrap {
		// DNS transports require their own resolver even when route has a
		// default_domain_resolver. Reuse one dependency for every hostname.
		compiled[index]["domain_resolver"] = map[string]any{"server": bootstrapTag}
	}
	return compiled, nil
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
	if rule.Timeout != "" {
		switch rule.Action {
		case C.RuleActionTypeRoute, C.RuleActionTypeEvaluate, C.RuleActionTypeRouteOptions:
			compiled["timeout"] = rule.Timeout
		default:
			return nil, fmt.Errorf("dns rule action %q does not support timeout", rule.Action)
		}
	}
	applyOptionalString(compiled, "client_subnet", rule.ClientSubnet)
	return compiled, nil
}
