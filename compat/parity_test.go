package compat

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	support "github.com/0xddy/sing-box-acp-support/singboxconfig"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	sharedprovider "github.com/acp/node-agent/api/provider"
	legacy "github.com/acp/node-agent/pkg/singboxconfig"
)

func TestSupportedOutboundTypesParity(t *testing.T) {
	want := legacy.SupportedOutboundTypes()
	got := support.SupportedOutboundTypes()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("supported outbound types = %v, legacy = %v", got, want)
	}
}

func TestCompileSnapshotParity(t *testing.T) {
	for _, testCase := range snapshotCases(t) {
		t.Run(testCase.name, func(t *testing.T) {
			want, legacyErr := legacy.CompileSnapshot(testCase.snapshot)
			got, supportErr := support.CompileSnapshot(testCase.snapshot)
			compareErrors(t, supportErr, legacyErr)
			if !bytes.Equal(got, want) {
				t.Fatalf("compiled output differs\nsupport:\n%s\nlegacy:\n%s", got, want)
			}
		})
	}
}

func TestValidateOutboundParity(t *testing.T) {
	testCases := []struct {
		name     string
		outbound *acpv1.OutboundConfig
	}{
		{name: "nil"},
		{name: "direct", outbound: outbound("direct", "bound-direct", `{}`)},
		{name: "selector", outbound: outbound("selector", "selector", `{"outbounds":["bound-direct"],"default":"bound-direct"}`)},
		{name: "urltest", outbound: outbound("urltest", "urltest", `{"outbounds":["bound-direct"],"url":"https://example.com/generate_204","interval":"5m"}`)},
		{name: "shadowsocks", outbound: outbound("shadowsocks", "ss", `{"server":"127.0.0.1","server_port":8388,"method":"aes-128-gcm","password":"secret"}`)},
		{name: "trojan", outbound: outbound("trojan", "trojan", `{"server":"127.0.0.1","server_port":443,"password":"secret"}`)},
		{name: "vless", outbound: outbound("vless", "vless", `{"server":"127.0.0.1","server_port":443,"uuid":"11111111-1111-4111-8111-111111111111"}`)},
		{name: "hysteria2", outbound: outbound("hysteria2", "hy2", `{"server":"127.0.0.1","server_port":443,"password":"secret","tls":{"enabled":true,"server_name":"example.com","insecure":true}}`)},
		{name: "unknown type", outbound: outbound("unknown", "unknown", `{}`)},
		{name: "unknown option", outbound: outbound("direct", "bad-direct", `{"unknown":true}`)},
		{name: "managed type", outbound: outbound("direct", "bad-direct", `{"type":"direct"}`)},
		{name: "array options", outbound: outbound("direct", "bad-direct", `[]`)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			compareErrors(t, support.ValidateOutbound(testCase.outbound), legacy.ValidateOutbound(testCase.outbound))
		})
	}
}

func TestDirectOutboundIsEmptyParity(t *testing.T) {
	for _, options := range []string{
		``,
		`{}`,
		`{"bind_interface":""}`,
		`{"inet4_bind_address":"192.0.2.10"}`,
		`{"domain_resolver":{"server":"default-dns"}}`,
		`{"udp_fragment":false}`,
		`{"unknown":true}`,
	} {
		t.Run(options, func(t *testing.T) {
			want, legacyErr := legacy.DirectOutboundIsEmpty([]byte(options))
			got, supportErr := support.DirectOutboundIsEmpty([]byte(options))
			compareErrors(t, supportErr, legacyErr)
			if got != want {
				t.Fatalf("empty = %v, legacy = %v", got, want)
			}
		})
	}
}

func snapshotCases(t *testing.T) []struct {
	name     string
	snapshot *acpv1.TopologySnapshot
} {
	t.Helper()
	vlessConfig := mustJSON(t, sharedprovider.VLESSRealityVisionConfig{
		Type: "vless", Tag: "node-vless", Listen: "::", ListenPort: 443,
		Flow: "xtls-rprx-vision", TCPFastOpen: true, Sniff: true,
		TLS: sharedprovider.VLESSRealityVisionTLSConfig{
			Enabled: true, ServerName: "www.example.com",
			Reality: sharedprovider.VLESSRealityConfig{
				Enabled: true, PrivateKey: "private-key", ShortID: []string{"0123456789abcdef"},
				Handshake: sharedprovider.RealityHandshake{Server: "www.example.com", ServerPort: 443},
			},
		},
	})
	hy2Config := mustJSON(t, sharedprovider.Hysteria2SalamanderConfig{
		Type: "hysteria2", Tag: "node-hy2", Listen: "::", ListenPort: 8443,
		UpMbps: 100, DownMbps: 200,
		Obfs: sharedprovider.Hysteria2ObfsConfig{Type: "salamander", Password: "obfs-secret"},
		TLS: sharedprovider.Hysteria2TLSConfig{
			Enabled: true, ServerName: "hy2.example.com", CertificatePEM: "certificate", PrivateKeyPEM: "private-key",
		},
	})

	return []struct {
		name     string
		snapshot *acpv1.TopologySnapshot
	}{
		{name: "nil"},
		{name: "vless route and dns", snapshot: &acpv1.TopologySnapshot{
			MachineId: "machine-1",
			Nodes: []*acpv1.NodeTopology{{
				NodeId: "node-vless", ProviderId: sharedprovider.VLESSRealityVisionID,
				ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfigJson: vlessConfig,
				Users: []*acpv1.UserCredential{{UserId: "user-1", Credential: "11111111-1111-4111-8111-111111111111"}},
			}},
			Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}},
			Route: &acpv1.RouteConfig{
				Final: "direct", DefaultDomainResolver: &acpv1.DomainResolveOptions{Server: "default-dns"},
			},
			Dns: &acpv1.DNSConfig{
				Servers: []*acpv1.DNSServer{{Type: "https", Tag: "default-dns", Server: "1.1.1.1"}}, Final: "default-dns",
			},
		}},
		{name: "hysteria2", snapshot: &acpv1.TopologySnapshot{
			MachineId: "machine-1",
			Nodes: []*acpv1.NodeTopology{{
				NodeId: "node-hy2", ProviderId: sharedprovider.Hysteria2SalamanderID,
				ProviderConfigVersion: sharedprovider.CurrentConfigVersion, ProviderConfigJson: hy2Config,
				Users: []*acpv1.UserCredential{{UserId: "user-7", Credential: "password"}},
			}},
			Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}},
		}},
		{name: "invalid route reference", snapshot: &acpv1.TopologySnapshot{
			MachineId: "machine-1",
			Outbounds: []*acpv1.OutboundConfig{{Type: "direct", Tag: "direct"}},
			Route:     &acpv1.RouteConfig{Final: "missing"},
		}},
	}
}

func outbound(outboundType, tag, options string) *acpv1.OutboundConfig {
	return &acpv1.OutboundConfig{Type: outboundType, Tag: tag, OptionsJson: []byte(options)}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func compareErrors(t *testing.T, got, want error) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Fatalf("support error = %v, legacy error = %v", got, want)
	}
	if got != nil && got.Error() != want.Error() {
		t.Fatalf("support error = %q, legacy error = %q", got, want)
	}
}
