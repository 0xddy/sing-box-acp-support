package compiler

import (
	"encoding/json"
	"testing"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
	sharedprovider "github.com/acp/node-agent/api/provider"
)

func providerJSON(t *testing.T, value any) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func vlessNode(t *testing.T, credential string) topology.NodeInstance {
	return topology.NodeInstance{NodeID: "node-vless-1", ProviderID: sharedprovider.VLESSRealityVisionID, ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfig: providerJSON(t, sharedprovider.VLESSRealityVisionConfig{Type: "vless", Tag: "node-vless-1", Listen: "::", ListenPort: 443, Flow: "xtls-rprx-vision", TCPFastOpen: true, Sniff: true, TLS: sharedprovider.VLESSRealityVisionTLSConfig{Enabled: true, ServerName: "www.example.com", Reality: sharedprovider.VLESSRealityConfig{Enabled: true, PrivateKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", ShortID: []string{"0123456789abcdef"}, Handshake: sharedprovider.RealityHandshake{Server: "www.example.com", ServerPort: 443}}}}), Users: []topology.UserCredential{{UserID: "1", Credential: credential}}}
}

func TestCompileVLESSRealityVisionProvider(t *testing.T) {
	cfg, err := Compile(topology.MachineTopology{MachineID: "machine-1", Nodes: []topology.NodeInstance{vlessNode(t, "11111111-1111-4111-8111-111111111111")}, Outbounds: []topology.Outbound{{Type: "direct", Tag: "direct"}}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err = json.Unmarshal(cfg, &decoded); err != nil {
		t.Fatal(err)
	}
	inbound := decoded["inbounds"].([]any)[0].(map[string]any)
	if inbound["tcp_fast_open"] != true {
		t.Fatalf("inbound=%+v", inbound)
	}
	multiplex, ok := inbound["multiplex"].(map[string]any)
	if !ok || multiplex["enabled"] != true {
		t.Fatalf("multiplex=%+v", inbound["multiplex"])
	}
	user := inbound["users"].([]any)[0].(map[string]any)
	if user["flow"] != "xtls-rprx-vision" {
		t.Fatalf("user=%+v", user)
	}
	rules := decoded["route"].(map[string]any)["rules"].([]any)
	if rules[0].(map[string]any)["action"] != "sniff" || rules[0].(map[string]any)["timeout"] != "300ms" {
		t.Fatalf("rules=%+v", rules)
	}
}

func TestCompileHysteria2SalamanderProvider(t *testing.T) {
	node := topology.NodeInstance{NodeID: "node-hy2", ProviderID: sharedprovider.Hysteria2SalamanderID, ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfig: providerJSON(t, sharedprovider.Hysteria2SalamanderConfig{Type: "hysteria2", Tag: "node-hy2", Listen: "::", ListenPort: 8443, UpMbps: 100, DownMbps: 200, Obfs: sharedprovider.Hysteria2ObfsConfig{Type: "salamander", Password: "obfs-secret"}, TLS: sharedprovider.Hysteria2TLSConfig{Enabled: true, ServerName: "hy2.example.com", CertificatePEM: "certificate", PrivateKeyPEM: "private-key"}}), Users: []topology.UserCredential{{UserID: "7", Credential: "user-password"}}}
	cfg, err := Compile(topology.MachineTopology{MachineID: "machine-1", Nodes: []topology.NodeInstance{node}, Outbounds: []topology.Outbound{{Type: "direct", Tag: "direct"}}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	_ = json.Unmarshal(cfg, &decoded)
	inbound := decoded["inbounds"].([]any)[0].(map[string]any)
	if inbound["type"] != "hysteria2" || inbound["obfs"].(map[string]any)["password"] != "obfs-secret" {
		t.Fatalf("inbound=%+v", inbound)
	}
	if inbound["users"].([]any)[0].(map[string]any)["password"] != "user-password" {
		t.Fatalf("users=%+v", inbound["users"])
	}
	rules := decoded["route"].(map[string]any)["rules"].([]any)
	rule := rules[0].(map[string]any)
	if len(rules) != 1 || rule["action"] != "sniff" || rule["timeout"] != "300ms" || rule["inbound"].([]any)[0] != "node-hy2" {
		t.Fatalf("Hysteria2 default sniff policy drifted: %+v", rules)
	}
}

func TestCompileHysteria2ProxyMasqueradeWithoutObfs(t *testing.T) {
	node := topology.NodeInstance{NodeID: "node-hy2", ProviderID: sharedprovider.Hysteria2SalamanderID, ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfig: providerJSON(t, sharedprovider.Hysteria2SalamanderConfig{Type: "hysteria2", Tag: "node-hy2", Listen: "::", ListenPort: 8443, Masquerade: &sharedprovider.Hysteria2MasqueradeConfig{Type: "proxy", URL: "https://www.example.com/landing", RewriteHost: true}, TLS: sharedprovider.Hysteria2TLSConfig{Enabled: true, ServerName: "hy2.example.com", CertificatePEM: "certificate", PrivateKeyPEM: "private-key"}}), Users: []topology.UserCredential{{UserID: "7", Credential: "user-password"}}}
	cfg, err := Compile(topology.MachineTopology{MachineID: "machine-1", Nodes: []topology.NodeInstance{node}, Outbounds: []topology.Outbound{{Type: "direct", Tag: "direct"}}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err = json.Unmarshal(cfg, &decoded); err != nil {
		t.Fatal(err)
	}
	inbound := decoded["inbounds"].([]any)[0].(map[string]any)
	if _, exists := inbound["obfs"]; exists {
		t.Fatalf("obfs should be omitted: %+v", inbound)
	}
	masquerade := inbound["masquerade"].(map[string]any)
	if masquerade["type"] != "proxy" || masquerade["url"] != "https://www.example.com/landing" || masquerade["rewrite_host"] != true {
		t.Fatalf("masquerade=%+v", masquerade)
	}
}

func TestCompileHysteria2FixedResponseMasqueradeWithoutObfs(t *testing.T) {
	node := topology.NodeInstance{NodeID: "node-hy2", ProviderID: sharedprovider.Hysteria2SalamanderID, ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfig: providerJSON(t, sharedprovider.Hysteria2SalamanderConfig{Type: "hysteria2", Tag: "node-hy2", Listen: "::", ListenPort: 8443, Masquerade: &sharedprovider.Hysteria2MasqueradeConfig{Type: "string", Content: "<h1>Welcome</h1>"}, TLS: sharedprovider.Hysteria2TLSConfig{Enabled: true, ServerName: "hy2.example.com", CertificatePEM: "certificate", PrivateKeyPEM: "private-key"}}), Users: []topology.UserCredential{{UserID: "7", Credential: "user-password"}}}
	cfg, err := Compile(topology.MachineTopology{MachineID: "machine-1", Nodes: []topology.NodeInstance{node}, Outbounds: []topology.Outbound{{Type: "direct", Tag: "direct"}}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err = json.Unmarshal(cfg, &decoded); err != nil {
		t.Fatal(err)
	}
	inbound := decoded["inbounds"].([]any)[0].(map[string]any)
	masquerade := inbound["masquerade"].(map[string]any)
	headers := masquerade["headers"].(map[string]any)
	if masquerade["type"] != "string" || masquerade["content"] != "<h1>Welcome</h1>" || headers["Content-Type"] != "text/html; charset=utf-8" {
		t.Fatalf("masquerade=%+v", masquerade)
	}
}

func TestCompileRejectsUnknownProvider(t *testing.T) {
	_, err := Compile(topology.MachineTopology{MachineID: "machine-1", Nodes: []topology.NodeInstance{{NodeID: "node-1", ProviderID: "unknown", ProviderConfigVersion: 1, ProviderConfig: json.RawMessage(`{}`)}}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompileIncludesRuleSetCacheFile(t *testing.T) {
	cfg, err := Compile(topology.MachineTopology{MachineID: "machine-1", Outbounds: []topology.Outbound{{Type: "direct", Tag: "direct"}}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	_ = json.Unmarshal(cfg, &decoded)
	cache := decoded["experimental"].(map[string]any)["cache_file"].(map[string]any)
	if cache["enabled"] != true {
		t.Fatalf("cache=%+v", cache)
	}
}
