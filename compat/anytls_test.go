package compat

import (
	"testing"

	support "github.com/0xddy/sing-box-acp-support/singboxconfig"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	sharedprovider "github.com/acp/node-agent/api/provider"
)

func TestAnyTLSCompilationParity(t *testing.T) {
	config := sharedprovider.AnyTLSTLSConfig{Type: "anytls", Tag: "node-anytls", Listen: "::", ListenPort: 443, TLS: sharedprovider.AnyTLSTLSSettings{Enabled: true, ServerName: "example.com", CertificatePEM: "certificate", PrivateKeyPEM: "private-key"}}
	for _, tc := range []struct {
		name    string
		web     sharedprovider.WebConfig
		suffix  string
		invalid bool
	}{
		{"static", sharedprovider.WebConfig{Type: "static", Title: "Website", Body: "About"}, "", false},
		{"proxy", sharedprovider.WebConfig{Type: "proxy", URL: "https://example.com/site"}, "", false},
		{"unknown-mode", sharedprovider.WebConfig{Type: "unknown"}, "", true},
		{"trailing-json", sharedprovider.WebConfig{Type: "static", Title: "Website"}, " {}", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config.Web = tc.web
			encoded := append(mustJSON(t, config), tc.suffix...)
			snapshot := &acpv1.TopologySnapshot{MachineId: "machine-1", Nodes: []*acpv1.NodeTopology{{NodeId: "node-anytls", ProviderId: sharedprovider.AnyTLSTLSID, ProviderConfigVersion: 1, ProviderConfigJson: encoded, Users: []*acpv1.UserCredential{{UserId: "7", Credential: "password"}}}}, Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}}}
			compareCompiledSnapshot(t, snapshot)
			if _, err := support.CompileSnapshot(snapshot); (err != nil) != tc.invalid {
				t.Fatalf("unexpected validation result: %v", err)
			}
		})
	}
}
