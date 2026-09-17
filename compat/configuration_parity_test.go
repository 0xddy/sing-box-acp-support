package compat

import (
	"bytes"
	"testing"

	support "github.com/0xddy/sing-box-acp-support/singboxconfig"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	sharedprovider "github.com/acp/node-agent/api/provider"
	legacy "github.com/acp/node-agent/pkg/singboxconfig"
)

func compareCompiledSnapshot(t *testing.T, snapshot *acpv1.TopologySnapshot) {
	t.Helper()
	want, legacyErr := legacy.CompileSnapshot(snapshot)
	got, supportErr := support.CompileSnapshot(snapshot)
	compareErrors(t, supportErr, legacyErr)
	if !bytes.Equal(got, want) {
		t.Fatalf("compiled output differs\nsupport:\n%s\nnode-agent:\n%s", got, want)
	}
}

func TestDNSBootstrapParity(t *testing.T) {
	for _, tc := range []struct {
		name, final string
		servers     []*acpv1.DNSServer
	}{
		{"prefer explicit literal final", "final", []*acpv1.DNSServer{
			{Type: "udp", Tag: "other-ip", Server: "8.8.8.8"},
			{Type: "https", Tag: "host", Server: "dns.google"},
			{Type: "https", Tag: "final", Server: "1.1.1.1"},
		}},
		{"hostname final reuses IPv6 resolver", "host", []*acpv1.DNSServer{
			{Type: "tls", Tag: "host", Server: "dns.google"},
			{Type: "https", Tag: "ipv6", Server: "2606:4700:4700::1111"},
		}},
		{"hostname transports share encrypted bootstrap", "https", []*acpv1.DNSServer{
			{Type: "https", Tag: "https", Server: "dns.google"},
			{Type: "quic", Tag: "quic", Server: "dns.example"},
			{Type: "tcp", Tag: "tcp", Server: "dns.example"},
			{Type: "udp", Tag: "udp", Server: "dns.example"},
		}},
		{"bootstrap tag collision", "__acp_dns_bootstrap", []*acpv1.DNSServer{
			{Type: "https", Tag: "__acp_dns_bootstrap", Server: "dns.google"},
			{Type: "tls", Tag: "__acp_dns_bootstrap_1", Server: "one.one.one.one"},
		}},
		{"detoured literal cannot bootstrap other servers", "detoured", []*acpv1.DNSServer{
			{Type: "https", Tag: "host", Server: "dns.google"},
			{Type: "https", Tag: "detoured", Server: "1.1.1.1", Detour: "direct"},
		}},
		{"detour owns hostname resolution", "host", []*acpv1.DNSServer{
			{Type: "https", Tag: "host", Server: "dns.google", Detour: "direct"},
		}},
		{"literal has no bootstrap dependency", "literal", []*acpv1.DNSServer{
			{Type: "https", Tag: "literal", Server: "1.1.1.1"},
		}},
		{"invalid server preserves validation error", "invalid", []*acpv1.DNSServer{
			{Type: "https", Tag: "invalid"},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			compareCompiledSnapshot(t, &acpv1.TopologySnapshot{
				MachineId: "machine-1",
				Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}},
				Dns:       &acpv1.DNSConfig{Final: tc.final, Servers: tc.servers},
			})
		})
	}
}

func TestRemoteRuleSetDefaultHTTPClientParity(t *testing.T) {
	for _, tc := range []struct {
		name, final, detour string
		outbounds           []*acpv1.OutboundConfig
	}{
		{name: "implicit direct"},
		{name: "configured direct dialer", final: "direct", outbounds: []*acpv1.OutboundConfig{
			outbound("direct", "direct", `{"tcp_fast_open":true,"bind_interface":"eth0"}`),
		}},
		{name: "selected default", final: "selected", outbounds: []*acpv1.OutboundConfig{
			outbound("direct", "direct", `{}`),
			outbound("selector", "selected", `{"outbounds":["direct"],"default":"direct"}`),
		}},
		{name: "first outbound without explicit final", outbounds: []*acpv1.OutboundConfig{
			outbound("selector", "selected", `{"outbounds":["direct"],"default":"direct"}`),
			outbound("direct", "direct", `{}`),
		}},
		{name: "explicit rule-set detour", detour: "direct", outbounds: []*acpv1.OutboundConfig{
			outbound("direct", "direct", `{"tcp_fast_open":true}`),
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			compareCompiledSnapshot(t, &acpv1.TopologySnapshot{
				MachineId: "machine-1", Outbounds: tc.outbounds,
				Route: &acpv1.RouteConfig{Final: tc.final, RuleSets: []*acpv1.RouteRuleSet{{
					Type: "remote", Tag: "rules", Format: "binary", Url: "https://example.com/rules.srs", DownloadDetour: tc.detour,
				}}},
			})
		})
	}
}

func TestDNSRuleValidationPrecedenceParity(t *testing.T) {
	compareCompiledSnapshot(t, &acpv1.TopologySnapshot{
		MachineId: "machine-1",
		Dns: &acpv1.DNSConfig{Rules: []*acpv1.DNSRule{{
			Action: "respond", Timeout: "5s", RewriteTtl: "not-a-number",
		}}},
	})
}

func TestVLESSRealityProviderSafetyParity(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*sharedprovider.VLESSRealityVisionConfig)
	}{
		{"valid", func(*sharedprovider.VLESSRealityVisionConfig) {}},
		{"TLS disabled", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Enabled = false }},
		{"REALITY disabled", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.Enabled = false }},
		{"missing handshake host", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.Handshake.Server = "" }},
		{"missing handshake port", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.Handshake.ServerPort = 0 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := sharedprovider.VLESSRealityVisionConfig{
				Type: "vless", ListenPort: 443, Sniff: true,
				TLS: sharedprovider.VLESSRealityVisionTLSConfig{
					Enabled: true, ServerName: "www.example.com",
					Reality: sharedprovider.VLESSRealityConfig{
						Enabled: true, PrivateKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", ShortID: []string{"0123456789abcdef"},
						Handshake: sharedprovider.RealityHandshake{Server: "www.example.com", ServerPort: 443},
					},
				},
			}
			tc.mutate(&cfg)
			snapshot := &acpv1.TopologySnapshot{MachineId: "machine-1", Nodes: []*acpv1.NodeTopology{{
				NodeId: "node-1", ProviderId: sharedprovider.VLESSRealityVisionID,
				ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfigJson: mustJSON(t, cfg),
			}}}
			compareCompiledSnapshot(t, snapshot)
			_, err := support.CompileSnapshot(snapshot)
			if (err == nil) != (tc.name == "valid") {
				t.Fatalf("REALITY safety acceptance mismatch: %v", err)
			}
		})
	}
}
