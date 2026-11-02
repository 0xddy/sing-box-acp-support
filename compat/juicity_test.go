package compat

import (
	"crypto/x509"
	"encoding/pem"
	support "github.com/0xddy/sing-box-acp-support/singboxconfig"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"github.com/acp/node-agent/api/provider"
	"net/http/httptest"
	"testing"
)

func TestJuicityCompilationParity(t *testing.T) {
	server := httptest.NewTLSServer(nil)
	defer server.Close()
	pair := server.TLS.Certificates[0]
	key, err := x509.MarshalPKCS8PrivateKey(pair.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	cfg := provider.JuicityConfig{
		JuicityListenOptions: provider.JuicityListenOptions{Type: "juicity", Tag: "juicity-in", Listen: "::", ListenPort: 443, MaxConcurrentStreams: 30},
		TLS: provider.CertificateTLSSettings{
			Enabled: true, ServerName: "example.com",
			CertificatePEM: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: pair.Certificate[0]})),
			PrivateKeyPEM:  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key})),
		},
	}
	node := &acpv1.NodeTopology{NodeId: "n", ProviderId: provider.JuicityID, ProviderConfigVersion: 1, ProviderConfigJson: mustJSON(t, cfg)}
	snapshot := &acpv1.TopologySnapshot{MachineId: "m", Nodes: []*acpv1.NodeTopology{node}}
	active := &acpv1.UserCredential{UserId: "42", Credential: "00112233-4455-4677-8899-aabbccddeeff", Status: acpv1.UserStatus_USER_STATUS_ACTIVE}
	disabled := &acpv1.UserCredential{UserId: "27", Status: acpv1.UserStatus_USER_STATUS_DISABLED}
	for _, tc := range []struct {
		name    string
		users   []*acpv1.UserCredential
		invalid bool
	}{
		{"empty", nil, false},
		{"active and disabled", []*acpv1.UserCredential{active, disabled}, false},
		{"invalid numeric fallback", []*acpv1.UserCredential{{UserId: "27", Credential: "27", Status: acpv1.UserStatus_USER_STATUS_ACTIVE}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			node.Users = tc.users
			compareCompiledSnapshot(t, snapshot)
			if _, err := support.CompileSnapshot(snapshot); (err != nil) != tc.invalid {
				t.Fatalf("validation: %v", err)
			}
		})
	}
}
