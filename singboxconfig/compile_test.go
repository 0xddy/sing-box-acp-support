package singboxconfig

import (
	"encoding/json"
	"testing"

	acpv1 "github.com/acp/node-agent/api/acp/v1"
	sharedprovider "github.com/acp/node-agent/api/provider"
)

func TestCompileSnapshotKeepsCredentialsAndRealityPrivateKey(t *testing.T) {
	providerConfig, _ := json.Marshal(sharedprovider.VLESSRealityVisionConfig{Type: "vless", Tag: "node-1", Listen: "::", ListenPort: 443, Flow: "xtls-rprx-vision", TLS: sharedprovider.VLESSRealityVisionTLSConfig{Enabled: true, ServerName: "www.example.com", Reality: sharedprovider.VLESSRealityConfig{Enabled: true, PrivateKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", ShortID: []string{"0123456789abcdef"}, Handshake: sharedprovider.RealityHandshake{Server: "www.example.com", ServerPort: 443}}}})
	compiled, err := CompileSnapshot(&acpv1.TopologySnapshot{MachineId: "machine-1", Nodes: []*acpv1.NodeTopology{{NodeId: "node-1", ProviderId: sharedprovider.VLESSRealityVisionID, ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfigJson: providerConfig, Users: []*acpv1.UserCredential{{UserId: "1", Credential: "11111111-1111-4111-8111-111111111111"}}}}, Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}}, Route: &acpv1.RouteConfig{Final: "direct"}})
	if err != nil {
		t.Fatalf("CompileSnapshot() error = %v", err)
	}
	var cfg map[string]any
	if err = json.Unmarshal(compiled, &cfg); err != nil {
		t.Fatal(err)
	}
	inbound := cfg["inbounds"].([]any)[0].(map[string]any)
	if got := inbound["users"].([]any)[0].(map[string]any)["uuid"]; got != "11111111-1111-4111-8111-111111111111" {
		t.Fatalf("uuid = %v", got)
	}
	if got := inbound["tls"].(map[string]any)["reality"].(map[string]any)["private_key"]; got != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" {
		t.Fatalf("private_key = %v", got)
	}
}

func TestCompileSnapshotEmitsDefaultDomainResolver(t *testing.T) {
	compiled, err := CompileSnapshot(&acpv1.TopologySnapshot{
		MachineId: "machine-1",
		Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}},
		Route: &acpv1.RouteConfig{
			Final: "direct",
			DefaultDomainResolver: &acpv1.DomainResolveOptions{
				Server: "default-dns",
			},
		},
		Dns: &acpv1.DNSConfig{
			Servers: []*acpv1.DNSServer{{
				Type:   "https",
				Tag:    "default-dns",
				Server: "1.1.1.1",
			}},
			Final: "default-dns",
		},
	})
	if err != nil {
		t.Fatalf("CompileSnapshot() error = %v", err)
	}
	var config map[string]any
	if err = json.Unmarshal(compiled, &config); err != nil {
		t.Fatal(err)
	}
	route := config["route"].(map[string]any)
	resolver := route["default_domain_resolver"].(map[string]any)
	if resolver["server"] != "default-dns" {
		t.Fatalf("route.default_domain_resolver.server = %v, want default-dns", resolver["server"])
	}
	dns := config["dns"].(map[string]any)
	if dns["final"] != "default-dns" {
		t.Fatalf("dns.final = %v, want default-dns", dns["final"])
	}
}

func TestCompileSnapshotRejectsInvalidRouteConversion(t *testing.T) {
	compiled, err := CompileSnapshot(&acpv1.TopologySnapshot{
		MachineId: "machine-1",
		Route: &acpv1.RouteConfig{Rules: []*acpv1.RouteRule{{
			Action: "reject", IpVersion: 300,
		}}},
	})
	if err == nil || compiled != nil {
		t.Fatalf("invalid route must not compile into an unrestricted reject: config=%s err=%v", compiled, err)
	}
}

func TestValidateOutboundUsesPinnedOptionTypes(t *testing.T) {
	if err := ValidateOutbound(&acpv1.OutboundConfig{
		Type:        "direct",
		Tag:         "bound-direct",
		OptionsJson: []byte(`{"inet4_bind_address":"203.0.113.10"}`),
	}); err != nil {
		t.Fatalf("validate Direct outbound: %v", err)
	}

	err := ValidateOutbound(&acpv1.OutboundConfig{
		Type:        "direct",
		Tag:         "invalid-direct",
		OptionsJson: []byte(`{"unknown":true}`),
	})
	if err == nil {
		t.Fatal("expected unknown Direct option to be rejected")
	}
}

func TestValidateOutboundRejectsNil(t *testing.T) {
	if err := ValidateOutbound(nil); err == nil {
		t.Fatal("expected nil outbound to be rejected")
	}
}

func TestValidateOutboundRejectsNestedUnknownFieldsAndInvalidDurations(t *testing.T) {
	testCases := []*acpv1.OutboundConfig{
		{
			Type:        "vless",
			Tag:         "invalid-vless",
			OptionsJson: []byte(`{"tls":{"enabled":true,"unknown":true}}`),
		},
		{
			Type:        "direct",
			Tag:         "invalid-direct",
			OptionsJson: []byte(`{"connect_timeout":"not-a-duration"}`),
		},
	}
	for _, outbound := range testCases {
		if err := ValidateOutbound(outbound); err == nil {
			t.Fatalf("expected %s options to be rejected", outbound.GetType())
		}
	}
}

func TestValidateOutboundLeavesRuntimeRequiredValuesToNode(t *testing.T) {
	// This package mirrors sing-box option decoding. Constructor/start-time
	// required-value checks remain a node runtime concern.
	for _, outboundType := range []string{"selector", "urltest", "shadowsocks", "trojan", "vless", "hysteria2"} {
		if err := ValidateOutbound(&acpv1.OutboundConfig{
			Type: outboundType, Tag: "zero-value-" + outboundType, OptionsJson: []byte(`{}`),
		}); err != nil {
			t.Fatalf("schema-valid zero-value %s options: %v", outboundType, err)
		}
	}
}

func TestSupportedOutboundTypesUsePinnedOptionTypes(t *testing.T) {
	want := []string{"direct", "selector", "urltest", "shadowsocks", "trojan", "vless", "hysteria2"}
	got := SupportedOutboundTypes()
	if len(got) != len(want) {
		t.Fatalf("supported type count = %d, want %d", len(got), len(want))
	}
	for index, outboundType := range want {
		if got[index] != outboundType {
			t.Fatalf("supported type[%d] = %q, want %q", index, got[index], outboundType)
		}
		if err := ValidateOutbound(&acpv1.OutboundConfig{
			Type: outboundType, Tag: "registry-test", OptionsJson: minimalOutboundOptions(outboundType),
		}); err != nil {
			t.Fatalf("registered type %q is not accepted: %v", outboundType, err)
		}
	}
}

func TestSupportedOutboundTypesReturnsCopy(t *testing.T) {
	first := SupportedOutboundTypes()
	first[0] = "modified"
	second := SupportedOutboundTypes()
	if second[0] != "direct" {
		t.Fatalf("supported type list was mutated: %v", second)
	}
}

func TestCompileSnapshotRejectsNil(t *testing.T) {
	if _, err := CompileSnapshot(nil); err == nil {
		t.Fatal("expected nil topology snapshot to be rejected")
	}
}

func TestCompileSnapshotSortsNodesAndFiltersDisabledUsers(t *testing.T) {
	snapshot := &acpv1.TopologySnapshot{
		MachineId: "machine-1",
		Nodes: []*acpv1.NodeTopology{
			vlessSnapshotNode(t, "node-b", []*acpv1.UserCredential{{
				UserId: "user-b", Credential: "22222222-2222-4222-8222-222222222222",
			}}),
			vlessSnapshotNode(t, "node-a", []*acpv1.UserCredential{
				{UserId: "disabled", Credential: "33333333-3333-4333-8333-333333333333", Status: acpv1.UserStatus_USER_STATUS_DISABLED},
				{UserId: "active", Credential: "11111111-1111-4111-8111-111111111111"},
			}),
		},
	}

	compiled, err := CompileSnapshot(snapshot)
	if err != nil {
		t.Fatalf("CompileSnapshot() error = %v", err)
	}
	var config map[string]any
	if err := json.Unmarshal(compiled, &config); err != nil {
		t.Fatal(err)
	}
	inbounds := config["inbounds"].([]any)
	if got := inbounds[0].(map[string]any)["tag"]; got != "node-a" {
		t.Fatalf("first inbound tag = %v, want node-a", got)
	}
	if got := inbounds[1].(map[string]any)["tag"]; got != "node-b" {
		t.Fatalf("second inbound tag = %v, want node-b", got)
	}
	users := inbounds[0].(map[string]any)["users"].([]any)
	if len(users) != 1 || users[0].(map[string]any)["name"] != "active" {
		t.Fatalf("node-a users = %#v, want only active user", users)
	}
}

func vlessSnapshotNode(t *testing.T, nodeID string, users []*acpv1.UserCredential) *acpv1.NodeTopology {
	t.Helper()
	providerConfig, err := json.Marshal(sharedprovider.VLESSRealityVisionConfig{
		Type: "vless", Tag: nodeID, Listen: "::", ListenPort: 443, Flow: "xtls-rprx-vision",
		TLS: sharedprovider.VLESSRealityVisionTLSConfig{
			Enabled: true, ServerName: "www.example.com",
			Reality: sharedprovider.VLESSRealityConfig{
				Enabled: true, PrivateKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", ShortID: []string{"0123456789abcdef"},
				Handshake: sharedprovider.RealityHandshake{Server: "www.example.com", ServerPort: 443},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return &acpv1.NodeTopology{
		NodeId: nodeID, ProviderId: sharedprovider.VLESSRealityVisionID,
		ProviderConfigVersion: sharedprovider.CurrentConfigVersion,
		ProviderConfigJson:    providerConfig,
		Users:                 users,
	}
}

func minimalOutboundOptions(outboundType string) []byte {
	switch outboundType {
	case "direct":
		return []byte(`{}`)
	case "selector":
		return []byte(`{"outbounds":["direct"]}`)
	case "urltest":
		return []byte(`{"outbounds":["direct"]}`)
	case "shadowsocks":
		return []byte(`{"server":"127.0.0.1","server_port":8388,"method":"aes-128-gcm","password":"secret"}`)
	case "trojan":
		return []byte(`{"server":"127.0.0.1","server_port":443,"password":"secret"}`)
	case "vless":
		return []byte(`{"server":"127.0.0.1","server_port":443,"uuid":"11111111-1111-4111-8111-111111111111"}`)
	case "hysteria2":
		return []byte(`{"server":"127.0.0.1","server_port":443,"password":"secret"}`)
	default:
		return []byte(`{}`)
	}
}
