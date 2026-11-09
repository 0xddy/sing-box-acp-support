package compat

import (
	"crypto/x509"
	"encoding/pem"
	"net/http/httptest"
	"testing"

	support "github.com/0xddy/sing-box-acp-support/singboxconfig"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"github.com/acp/node-agent/api/provider"
)

func TestNaiveCompilationParity(t *testing.T) {
	server := httptest.NewTLSServer(nil)
	defer server.Close()
	pair := server.TLS.Certificates[0]
	key, err := x509.MarshalPKCS8PrivateKey(pair.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	cfg := provider.NaiveTLSConfig{
		Type: "naive", Tag: "naive-in", Listen: "::", ListenPort: 443,
		Network: []string{"tcp", "udp"}, SubscriptionTransport: "https",
		TLS: provider.NaiveTLSSettings{
			Enabled: true, ServerName: "example.com", TrustMode: "system",
			CertificatePEM: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: pair.Certificate[0]})),
			PrivateKeyPEM:  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key})),
		},
		Web: provider.WebConfig{Type: "proxy", URL: "https://example.com"},
	}
	active := &acpv1.UserCredential{UserId: "42", Name: "display name", Credential: "password", Status: acpv1.UserStatus_USER_STATUS_ACTIVE}
	disabled := &acpv1.UserCredential{UserId: "27", Status: acpv1.UserStatus_USER_STATUS_DISABLED}
	for _, tc := range []struct {
		name    string
		users   []*acpv1.UserCredential
		web     provider.WebConfig
		invalid bool
	}{
		{"preview", nil, cfg.Web, false},
		{"active and disabled", []*acpv1.UserCredential{active, disabled}, cfg.Web, false},
		{"duplicate identity", []*acpv1.UserCredential{active, active}, cfg.Web, true},
		{"invalid web", nil, provider.WebConfig{Type: "unknown"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := cfg
			config.Web = tc.web
			snapshot := &acpv1.TopologySnapshot{MachineId: "m", Nodes: []*acpv1.NodeTopology{{NodeId: "naive-in", ProviderId: provider.NaiveTLSID, ProviderConfigVersion: 1, ProviderConfigJson: mustJSON(t, config), Users: tc.users}}}
			compareCompiledSnapshot(t, snapshot)
			if _, err := support.CompileSnapshot(snapshot); (err != nil) != tc.invalid {
				t.Fatalf("validation: %v", err)
			}
		})
	}
}
