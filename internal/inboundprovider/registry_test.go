package inboundprovider

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
	sharedprovider "github.com/acp/node-agent/api/provider"
)

func hysteria2Node(t *testing.T, mutate func(*sharedprovider.Hysteria2SalamanderConfig)) topology.NodeInstance {
	t.Helper()
	cfg := sharedprovider.Hysteria2SalamanderConfig{
		Type:       "hysteria2",
		Tag:        "node-hy2",
		Listen:     "::",
		ListenPort: 8443,
		Obfs:       sharedprovider.Hysteria2ObfsConfig{Type: "salamander", Password: "obfs-secret"},
		TLS: sharedprovider.Hysteria2TLSConfig{
			Enabled:        true,
			ServerName:     "hy2.example.com",
			CertificatePEM: "certificate",
			PrivateKeyPEM:  "private-key",
		},
	}
	if mutate != nil {
		mutate(&cfg)
	}
	encoded, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal provider config: %v", err)
	}
	return topology.NodeInstance{
		NodeID:                "node-1",
		ProviderID:            sharedprovider.Hysteria2SalamanderID,
		ProviderConfigVersion: sharedprovider.CurrentConfigVersion,
		ProviderConfig:        encoded,
	}
}

func testUsers() []topology.UserCredential {
	return []topology.UserCredential{{UserID: "user-1", Credential: "secret"}}
}

// Each direction maps to an independent sing-box bandwidth setting, so a node
// that only caps one direction must keep that cap.
func TestHysteria2BandwidthDirectionsAreIndependent(t *testing.T) {
	testCases := []struct {
		name         string
		upMbps       int
		downMbps     int
		wantUpMbps   any
		wantDownMbps any
	}{
		{name: "both", upMbps: 100, downMbps: 200, wantUpMbps: 100, wantDownMbps: 200},
		{name: "download only", downMbps: 200, wantDownMbps: 200},
		{name: "upload only", upMbps: 100, wantUpMbps: 100},
		{name: "neither"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			node := hysteria2Node(t, func(cfg *sharedprovider.Hysteria2SalamanderConfig) {
				cfg.UpMbps = testCase.upMbps
				cfg.DownMbps = testCase.downMbps
			})
			result, err := DefaultRegistry().Build(node, testUsers())
			if err != nil {
				t.Fatalf("build hysteria2 inbound: %v", err)
			}
			if got := result.Inbound["up_mbps"]; got != testCase.wantUpMbps {
				t.Fatalf("up_mbps = %v, want %v", got, testCase.wantUpMbps)
			}
			if got := result.Inbound["down_mbps"]; got != testCase.wantDownMbps {
				t.Fatalf("down_mbps = %v, want %v", got, testCase.wantDownMbps)
			}
		})
	}
}

func TestHysteria2RejectsNegativeBandwidth(t *testing.T) {
	node := hysteria2Node(t, func(cfg *sharedprovider.Hysteria2SalamanderConfig) {
		cfg.DownMbps = -1
	})
	_, err := DefaultRegistry().Build(node, testUsers())
	if err == nil || !strings.Contains(err.Error(), "must not be negative") {
		t.Fatalf("build error = %v, want a negative bandwidth rejection", err)
	}
}
