package compiler

import (
	"errors"
	"fmt"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func compileRoute(route *topology.Route, outboundTags map[string]struct{}, defaultFinal string) (map[string]any, error) {
	if route == nil {
		if defaultFinal == "" {
			return nil, nil
		}
		return map[string]any{"final": defaultFinal}, nil
	}

	compiled := map[string]any{}

	if len(route.Rules) > 0 {
		rules := make([]map[string]any, 0, len(route.Rules))
		for _, rule := range route.Rules {
			compiledRule, err := compileRouteRule(rule, outboundTags)
			if err != nil {
				return nil, err
			}
			rules = append(rules, compiledRule)
		}
		compiled["rules"] = rules
	}
	if len(route.RuleSets) > 0 {
		ruleSets := make([]map[string]any, 0, len(route.RuleSets))
		for _, ruleSet := range route.RuleSets {
			compiledRuleSet, err := compileRouteRuleSet(ruleSet)
			if err != nil {
				return nil, err
			}
			ruleSets = append(ruleSets, compiledRuleSet)
		}
		compiled["rule_set"] = ruleSets
	}
	if route.Final != "" {
		if _, ok := outboundTags[route.Final]; !ok {
			return nil, fmt.Errorf("route final references unknown outbound %q", route.Final)
		}
		compiled["final"] = route.Final
	}
	if route.AutoDetectInterface {
		compiled["auto_detect_interface"] = route.AutoDetectInterface
	}
	if route.DefaultInterface != "" {
		compiled["default_interface"] = route.DefaultInterface
	}
	if route.DefaultMark != 0 {
		compiled["default_mark"] = route.DefaultMark
	}
	if route.FindProcess {
		compiled["find_process"] = route.FindProcess
	}
	applyNestedStruct(compiled, "geoip", route.GeoIP)
	applyNestedStruct(compiled, "geosite", route.Geosite)
	if route.OverrideAndroidVPN {
		compiled["override_android_vpn"] = route.OverrideAndroidVPN
	}
	applyNestedStruct(compiled, "default_domain_resolver", route.DefaultDomainResolver)
	applyNestedStruct(compiled, "default_network_strategy", route.DefaultNetworkStrategy)
	applyStringList(compiled, "default_network_type", route.DefaultNetworkType)
	applyStringList(compiled, "default_fallback_network_type", route.DefaultFallbackNetworkType)
	applyOptionalString(compiled, "default_fallback_delay", route.DefaultFallbackDelay)
	return compiled, nil
}

func compileRouteRule(rule topology.RouteRule, outboundTags map[string]struct{}) (map[string]any, error) {
	compiled := map[string]any{}

	if rule.Type != "" {
		compiled["type"] = rule.Type
	}
	if len(rule.Inbound) > 0 {
		compiled["inbound"] = rule.Inbound
	}
	if len(rule.Network) > 0 {
		compiled["network"] = rule.Network
	}
	if rule.IPVersion != 0 {
		compiled["ip_version"] = rule.IPVersion
	}
	if len(rule.Domain) > 0 {
		compiled["domain"] = rule.Domain
	}
	if len(rule.DomainSuffix) > 0 {
		compiled["domain_suffix"] = rule.DomainSuffix
	}
	if len(rule.DomainKeyword) > 0 {
		compiled["domain_keyword"] = rule.DomainKeyword
	}
	if len(rule.DomainRegex) > 0 {
		compiled["domain_regex"] = rule.DomainRegex
	}
	if len(rule.SourceIPCidr) > 0 {
		compiled["source_ip_cidr"] = rule.SourceIPCidr
	}
	if len(rule.IPCidr) > 0 {
		compiled["ip_cidr"] = rule.IPCidr
	}
	if rule.SourceIPIsPrivate != nil {
		compiled["source_ip_is_private"] = *rule.SourceIPIsPrivate
	}
	if rule.IPIsPrivate != nil {
		compiled["ip_is_private"] = *rule.IPIsPrivate
	}
	if len(rule.Port) > 0 {
		compiled["port"] = rule.Port
	}
	if len(rule.PortRange) > 0 {
		compiled["port_range"] = rule.PortRange
	}
	if len(rule.SourcePortRange) > 0 {
		compiled["source_port_range"] = rule.SourcePortRange
	}
	if len(rule.Protocol) > 0 {
		compiled["protocol"] = rule.Protocol
	}
	applyStringList(compiled, "auth_user", rule.AuthUser)
	applyStringList(compiled, "client", rule.Client)
	applyStringList(compiled, "geosite", rule.Geosite)
	applyStringList(compiled, "source_geoip", rule.SourceGeoIP)
	applyStringList(compiled, "geoip", rule.GeoIP)
	applyUint32List(compiled, "source_port", rule.SourcePort)
	applyStringList(compiled, "process_name", rule.ProcessName)
	applyStringList(compiled, "process_path", rule.ProcessPath)
	applyStringList(compiled, "process_path_regex", rule.ProcessPathRegex)
	applyStringList(compiled, "package_name", rule.PackageName)
	applyStringList(compiled, "user", rule.User)
	applyInt32List(compiled, "user_id", rule.UserID)
	applyOptionalString(compiled, "clash_mode", rule.ClashMode)
	applyStringList(compiled, "network_type", rule.NetworkType)
	if rule.NetworkIsExpensive != nil {
		compiled["network_is_expensive"] = *rule.NetworkIsExpensive
	}
	if rule.NetworkIsConstrained != nil {
		compiled["network_is_constrained"] = *rule.NetworkIsConstrained
	}
	applyStringList(compiled, "wifi_ssid", rule.WIFISSID)
	applyStringList(compiled, "wifi_bssid", rule.WIFIBSSID)
	applyStringList(compiled, "default_interface_address", rule.DefaultInterfaceAddress)
	applyStringList(compiled, "preferred_by", rule.PreferredBy)
	if len(rule.RuleSet) > 0 {
		compiled["rule_set"] = rule.RuleSet
	}
	if rule.RuleSetIPCIDRMatchSource {
		compiled["rule_set_ip_cidr_match_source"] = rule.RuleSetIPCIDRMatchSource
	}
	if rule.Invert {
		compiled["invert"] = rule.Invert
	}
	if rule.Mode != "" {
		compiled["mode"] = rule.Mode
	}
	if len(rule.Rules) > 0 {
		rules := make([]map[string]any, 0, len(rule.Rules))
		for _, nested := range rule.Rules {
			compiledNested, err := compileRouteRule(nested, outboundTags)
			if err != nil {
				return nil, err
			}
			rules = append(rules, compiledNested)
		}
		compiled["rules"] = rules
	}
	if rule.Outbound != "" {
		if _, ok := outboundTags[rule.Outbound]; !ok {
			return nil, fmt.Errorf("route rule references unknown outbound %q", rule.Outbound)
		}
		if rule.Action == "" {
			compiled["action"] = "route"
		}
		compiled["outbound"] = rule.Outbound
	}
	if rule.Action != "" {
		compiled["action"] = rule.Action
	}
	if rule.RouteOptions != nil {
		mergeStruct(compiled, rule.RouteOptions)
	}
	if rule.DirectOptions != nil {
		applyDialerOptions(compiled, *rule.DirectOptions)
	}
	if rule.SniffOptions != nil {
		mergeStruct(compiled, rule.SniffOptions)
	}
	if rule.ResolveOptions != nil {
		mergeStruct(compiled, rule.ResolveOptions)
	}
	if rule.Method != "" {
		compiled["method"] = rule.Method
	}
	if rule.NoDrop {
		compiled["no_drop"] = rule.NoDrop
	}
	return compiled, nil
}

func compileRouteRuleSet(ruleSet topology.RouteRuleSet) (map[string]any, error) {
	if ruleSet.Type == "" {
		return nil, errors.New("route rule_set type is required")
	}
	if ruleSet.Tag == "" {
		return nil, fmt.Errorf("route rule_set %s tag is required", ruleSet.Type)
	}
	compiled := map[string]any{
		"type": ruleSet.Type,
		"tag":  ruleSet.Tag,
	}
	if ruleSet.Format != "" {
		compiled["format"] = ruleSet.Format
	}
	if ruleSet.Path != "" {
		compiled["path"] = ruleSet.Path
	}
	if ruleSet.URL != "" {
		compiled["url"] = ruleSet.URL
	}
	if ruleSet.DownloadDetour != "" {
		compiled["http_client"] = map[string]any{
			"detour": ruleSet.DownloadDetour,
		}
	}
	if ruleSet.UpdateInterval != "" {
		compiled["update_interval"] = ruleSet.UpdateInterval
	}
	if len(ruleSet.Rules) > 0 {
		rules := make([]map[string]any, 0, len(ruleSet.Rules))
		for _, rule := range ruleSet.Rules {
			rules = append(rules, compileHeadlessRule(rule))
		}
		compiled["rules"] = rules
	}
	return compiled, nil
}

func compileHeadlessRule(rule topology.HeadlessRule) map[string]any {
	compiled := map[string]any{}
	applyOptionalString(compiled, "type", rule.Type)
	applyStringList(compiled, "network", rule.Network)
	applyStringList(compiled, "domain", rule.Domain)
	applyStringList(compiled, "domain_suffix", rule.DomainSuffix)
	applyStringList(compiled, "domain_keyword", rule.DomainKeyword)
	applyStringList(compiled, "domain_regex", rule.DomainRegex)
	applyStringList(compiled, "source_ip_cidr", rule.SourceIPCidr)
	applyStringList(compiled, "ip_cidr", rule.IPCidr)
	applyUint32List(compiled, "source_port", rule.SourcePort)
	applyStringList(compiled, "source_port_range", rule.SourcePortRange)
	applyUint32List(compiled, "port", rule.Port)
	applyStringList(compiled, "port_range", rule.PortRange)
	applyStringList(compiled, "process_name", rule.ProcessName)
	applyStringList(compiled, "process_path", rule.ProcessPath)
	applyStringList(compiled, "process_path_regex", rule.ProcessPathRegex)
	applyStringList(compiled, "package_name", rule.PackageName)
	applyStringList(compiled, "network_type", rule.NetworkType)
	if rule.NetworkIsExpensive != nil {
		compiled["network_is_expensive"] = *rule.NetworkIsExpensive
	}
	if rule.NetworkIsConstrained != nil {
		compiled["network_is_constrained"] = *rule.NetworkIsConstrained
	}
	applyStringList(compiled, "wifi_ssid", rule.WIFISSID)
	applyStringList(compiled, "wifi_bssid", rule.WIFIBSSID)
	applyStringList(compiled, "default_interface_address", rule.DefaultInterfaceAddress)
	if rule.Invert {
		compiled["invert"] = rule.Invert
	}
	applyOptionalString(compiled, "mode", rule.Mode)
	if len(rule.Rules) > 0 {
		rules := make([]map[string]any, 0, len(rule.Rules))
		for _, nested := range rule.Rules {
			rules = append(rules, compileHeadlessRule(nested))
		}
		compiled["rules"] = rules
	}
	return compiled
}
