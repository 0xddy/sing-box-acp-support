package compiler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/0xddy/sing-box-acp-support/internal/outboundvalidate"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func TestCompileOutboundUsesRegisteredSingBoxOptions(t *testing.T) {
	testCases := []struct {
		name         string
		outboundType string
		options      string
	}{
		{name: "direct", outboundType: "direct", options: `{"inet4_bind_address":"127.0.0.1"}`},
		{name: "selector", outboundType: "selector", options: `{"outbounds":["direct"],"default":"direct"}`},
		{name: "urltest", outboundType: "urltest", options: `{"outbounds":["direct"],"url":"https://example.com/generate_204","interval":"5m"}`},
		{name: "shadowsocks", outboundType: "shadowsocks", options: `{"server":"127.0.0.1","server_port":8388,"method":"aes-128-gcm","password":"secret"}`},
		{name: "trojan", outboundType: "trojan", options: `{"server":"127.0.0.1","server_port":443,"password":"secret"}`},
		{name: "vless", outboundType: "vless", options: `{"server":"127.0.0.1","server_port":443,"uuid":"11111111-1111-4111-8111-111111111111"}`},
		{name: "hysteria2", outboundType: "hysteria2", options: `{"server":"127.0.0.1","server_port":443,"password":"secret","tls":{"enabled":true,"server_name":"example.com","insecure":true}}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			entry, err := compileOutbound(topology.Outbound{
				Type:    testCase.outboundType,
				Tag:     "test-" + testCase.name,
				Options: json.RawMessage(testCase.options),
			})
			if err != nil {
				t.Fatalf("compile outbound: %v", err)
			}
			if got := entry["type"]; got != testCase.outboundType {
				t.Fatalf("type = %#v, want %q", got, testCase.outboundType)
			}
			if got := entry["tag"]; got != "test-"+testCase.name {
				t.Fatalf("tag = %#v", got)
			}
		})
	}
}

func TestCompileOutboundRejectsInvalidOptions(t *testing.T) {
	testCases := []struct {
		name      string
		outbound  topology.Outbound
		wantError string
	}{
		{
			name:      "unknown type",
			outbound:  topology.Outbound{Type: "unknown", Tag: "test", Options: json.RawMessage(`{}`)},
			wantError: "unknown outbound type",
		},
		{
			name:      "unknown field",
			outbound:  topology.Outbound{Type: "direct", Tag: "test", Options: json.RawMessage(`{"unknown":true}`)},
			wantError: "unknown field",
		},
		{
			name:      "managed type",
			outbound:  topology.Outbound{Type: "direct", Tag: "test", Options: json.RawMessage(`{"type":"direct"}`)},
			wantError: "managed field",
		},
		{
			name:      "managed tag case insensitive",
			outbound:  topology.Outbound{Type: "direct", Tag: "test", Options: json.RawMessage(`{"Tag":"other"}`)},
			wantError: "managed field",
		},
		{
			name:      "array",
			outbound:  topology.Outbound{Type: "direct", Tag: "test", Options: json.RawMessage(`[]`)},
			wantError: "JSON object",
		},
		{
			name:      "null",
			outbound:  topology.Outbound{Type: "direct", Tag: "test", Options: json.RawMessage(`null`)},
			wantError: "JSON object",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := compileOutbound(testCase.outbound)
			if err == nil || !strings.Contains(err.Error(), testCase.wantError) {
				t.Fatalf("compile error = %v, want containing %q", err, testCase.wantError)
			}
		})
	}
}

func TestCompileOutboundPreservesRawOptions(t *testing.T) {
	entry, err := compileOutbound(topology.Outbound{
		Type: "direct",
		Tag:  "bound-direct",
		Options: json.RawMessage(`{
			"inet4_bind_address":"203.0.113.10",
			"domain_resolver":{"server":"dns-via-bound-direct","strategy":"ipv4_only"}
		}`),
	})
	if err != nil {
		t.Fatalf("compile outbound: %v", err)
	}
	encoded, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("marshal compiled outbound: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("decode compiled outbound: %v", err)
	}
	if got := decoded["inet4_bind_address"]; got != "203.0.113.10" {
		t.Fatalf("inet4_bind_address = %#v", got)
	}
	resolver, ok := decoded["domain_resolver"].(map[string]any)
	if !ok {
		t.Fatalf("domain_resolver = %#v", decoded["domain_resolver"])
	}
	if got := resolver["server"]; got != "dns-via-bound-direct" {
		t.Fatalf("domain_resolver.server = %#v", got)
	}
}

func TestCompileOutboundHysteria2ChromeParrotCompatibility(t *testing.T) {
	tests := []struct {
		name         string
		outboundType string
		options      string
		wantPresent  bool
		wantValue    bool
	}{
		{
			name:         "missing defaults disabled",
			outboundType: "hysteria2",
			options:      `{}`,
			wantPresent:  true,
			wantValue:    true,
		},
		{
			name:         "explicit false is preserved",
			outboundType: "hysteria2",
			options:      `{"disable_chrome_parrot":false}`,
			wantPresent:  true,
			wantValue:    false,
		},
		{
			name:         "explicit true is preserved",
			outboundType: "hysteria2",
			options:      `{"disable_chrome_parrot":true}`,
			wantPresent:  true,
			wantValue:    true,
		},
		{
			name:         "other outbound is unchanged",
			outboundType: "direct",
			options:      `{}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := compileOutbound(topology.Outbound{
				Type:    test.outboundType,
				Tag:     "test",
				Options: json.RawMessage(test.options),
			})
			if err != nil {
				t.Fatalf("compile outbound: %v", err)
			}

			encoded, err := json.Marshal(entry)
			if err != nil {
				t.Fatalf("marshal compiled outbound: %v", err)
			}
			var decoded map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("decode compiled outbound: %v", err)
			}
			rawValue, present := decoded["disable_chrome_parrot"]
			if present != test.wantPresent {
				t.Fatalf("disable_chrome_parrot present = %v, want %v; outbound = %s", present, test.wantPresent, encoded)
			}
			if !present {
				return
			}
			var value bool
			if err := json.Unmarshal(rawValue, &value); err != nil {
				t.Fatalf("decode disable_chrome_parrot: %v", err)
			}
			if value != test.wantValue {
				t.Fatalf("disable_chrome_parrot = %v, want %v", value, test.wantValue)
			}
		})
	}
}

func TestCompileOutboundPreservesCaseFoldedHysteria2ChromeParrotOption(t *testing.T) {
	entry, err := compileOutbound(topology.Outbound{
		Type:    "hysteria2",
		Tag:     "test",
		Options: json.RawMessage(`{"Disable_Chrome_Parrot":false}`),
	})
	if err != nil {
		t.Fatalf("compile outbound: %v", err)
	}
	if _, exists := entry["disable_chrome_parrot"]; exists {
		t.Fatalf("case-folded option was duplicated: %+v", entry)
	}
	rawValue, ok := entry["Disable_Chrome_Parrot"].(json.RawMessage)
	if !ok {
		t.Fatalf("case-folded option type = %T, want json.RawMessage", entry["Disable_Chrome_Parrot"])
	}
	var value bool
	if err := json.Unmarshal(rawValue, &value); err != nil {
		t.Fatalf("decode case-folded option: %v", err)
	}
	if value {
		t.Fatal("case-folded explicit false was not preserved")
	}
}

func TestDirectOutboundIsEmptyMatchesRuntimeDialerSemantics(t *testing.T) {
	tests := []struct {
		name    string
		options string
		empty   bool
	}{
		{name: "empty object", options: `{}`, empty: true},
		{name: "zero value field", options: `{"bind_interface":""}`, empty: true},
		{name: "source address", options: `{"inet4_bind_address":"192.0.2.10"}`},
		{name: "domain resolver", options: `{"domain_resolver":{"server":"default-dns"}}`},
		{name: "explicit UDP fragment", options: `{"udp_fragment":false}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			empty, err := outboundvalidate.DirectOutboundIsEmpty(json.RawMessage(test.options))
			if err != nil {
				t.Fatalf("check Direct options: %v", err)
			}
			if empty != test.empty {
				t.Fatalf("empty = %v, want %v", empty, test.empty)
			}
		})
	}
	if _, err := outboundvalidate.DirectOutboundIsEmpty(json.RawMessage(`{"unknown":true}`)); err == nil {
		t.Fatal("expected unknown Direct option to fail")
	}
}
