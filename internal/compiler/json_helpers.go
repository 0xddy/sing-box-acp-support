package compiler

import (
	"encoding/json"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func applyDialerOptions(target map[string]any, dialer topology.DialerOptions) {
	applyOptionalString(target, "detour", dialer.Detour)
	applyOptionalString(target, "bind_interface", dialer.BindInterface)
	applyOptionalString(target, "inet4_bind_address", dialer.Inet4BindAddress)
	applyOptionalString(target, "inet6_bind_address", dialer.Inet6BindAddress)
	if dialer.RoutingMark != 0 {
		target["routing_mark"] = dialer.RoutingMark
	}
	if dialer.ReuseAddr {
		target["reuse_addr"] = dialer.ReuseAddr
	}
	applyOptionalString(target, "connect_timeout", dialer.ConnectTimeout)
	if dialer.TCPFastOpen {
		target["tcp_fast_open"] = dialer.TCPFastOpen
	}
	if dialer.TCPMultiPath {
		target["tcp_multi_path"] = dialer.TCPMultiPath
	}
	if dialer.UDPFragment {
		target["udp_fragment"] = dialer.UDPFragment
	}
	applyOptionalString(target, "udp_timeout", dialer.UDPTimeout)
	applyOptionalString(target, "domain_strategy", dialer.DomainStrategy)
	if dialer.BindAddressNoPort {
		target["bind_address_no_port"] = dialer.BindAddressNoPort
	}
	applyOptionalString(target, "protect_path", dialer.ProtectPath)
	applyOptionalString(target, "netns", dialer.NetNS)
	if dialer.DisableTCPKeepAlive {
		target["disable_tcp_keep_alive"] = dialer.DisableTCPKeepAlive
	}
	applyOptionalString(target, "tcp_keep_alive", dialer.TCPKeepAlive)
	applyOptionalString(target, "tcp_keep_alive_interval", dialer.TCPKeepAliveInterval)
	applyNestedStruct(target, "domain_resolver", dialer.DomainResolver)
	applyNestedStruct(target, "network_strategy", dialer.NetworkStrategy)
	applyStringList(target, "network_type", dialer.NetworkType)
	applyStringList(target, "fallback_network_type", dialer.FallbackNetworkType)
	applyOptionalString(target, "fallback_delay", dialer.FallbackDelay)
}

func applyOptionalString(target map[string]any, key string, value string) {
	if value != "" {
		target[key] = value
	}
}

func applyStringList(target map[string]any, key string, values []string) {
	if len(values) > 0 {
		target[key] = values
	}
}

func applyUint32List(target map[string]any, key string, values []uint32) {
	if len(values) > 0 {
		target[key] = values
	}
}

func applyInt32List(target map[string]any, key string, values []int32) {
	if len(values) > 0 {
		target[key] = values
	}
}

func applyNestedStruct(target map[string]any, key string, value any) {
	compiled := structMap(value)
	if len(compiled) > 0 {
		target[key] = compiled
	}
}

func mergeStruct(target map[string]any, value any) {
	for key, fieldValue := range structMap(value) {
		target[key] = fieldValue
	}
}

func structMap(value any) map[string]any {
	if value == nil {
		return nil
	}
	encoded, err := json.Marshal(value)
	if err != nil || string(encoded) == "null" || string(encoded) == "{}" {
		return nil
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil
	}
	return decoded
}
