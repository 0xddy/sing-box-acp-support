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

func vlessNode(t *testing.T, mutate func(*sharedprovider.VLESSRealityVisionConfig)) topology.NodeInstance {
	t.Helper()
	cfg := sharedprovider.VLESSRealityVisionConfig{
		Type: "vless", ListenPort: 443, Flow: topology.VLESSFlowRealityVision,
		TLS: sharedprovider.VLESSRealityVisionTLSConfig{
			Enabled: true, ServerName: "www.example.com",
			Reality: sharedprovider.VLESSRealityConfig{
				Enabled: true, PrivateKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
				ShortID:   []string{"0123456789abcdef"},
				Handshake: sharedprovider.RealityHandshake{Server: "www.example.com", ServerPort: 443},
			},
		},
	}
	if mutate != nil {
		mutate(&cfg)
	}
	encoded, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return topology.NodeInstance{
		NodeID: "node-1", ProviderID: sharedprovider.VLESSRealityVisionID,
		ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfig: encoded,
	}
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

func TestVLESSRejectsMalformedRealityParameters(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*sharedprovider.VLESSRealityVisionConfig)
		field  string
	}{
		{"missing SNI", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.ServerName = "" }, "server_name"},
		{"invalid key encoding", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.PrivateKey = "bad-secret!" }, "private_key"},
		{"short key", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.PrivateKey = "c2hvcnQ" }, "private_key"},
		{"long short ID", func(c *sharedprovider.VLESSRealityVisionConfig) {
			c.TLS.Reality.ShortID = []string{"0123456789abcdef00"}
		}, "short_id[0]"},
		{"odd short ID", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.ShortID = []string{"123"} }, "short_id[0]"},
		{"non-hex short ID", func(c *sharedprovider.VLESSRealityVisionConfig) { c.TLS.Reality.ShortID = []string{"xyzsecret"} }, "short_id[0]"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := DefaultRegistry().Build(vlessNode(t, testCase.mutate), testUsers())
			if err == nil || !strings.Contains(err.Error(), testCase.field) {
				t.Fatalf("error=%v, want field %s", err, testCase.field)
			}
		})
	}
	for _, shortID := range []string{"", "00", "0123456789abcdef"} {
		if _, err := DefaultRegistry().Build(vlessNode(t, func(c *sharedprovider.VLESSRealityVisionConfig) {
			c.TLS.Reality.ShortID = []string{shortID}
		}), testUsers()); err != nil {
			t.Fatalf("valid short ID %q rejected: %v", shortID, err)
		}
	}
}
