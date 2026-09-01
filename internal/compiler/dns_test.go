package compiler

import (
	"testing"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func TestCompileDNSUsesEncryptedDefaultServer(t *testing.T) {
	for _, testCase := range []struct {
		name string
		dns  *topology.DNS
	}{
		{name: "nil config"},
		{name: "empty config", dns: &topology.DNS{}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			compiled, err := compileDNS(testCase.dns)
			if err != nil {
				t.Fatalf("compileDNS() error = %v", err)
			}

			servers, ok := compiled["servers"].([]map[string]any)
			if !ok {
				t.Fatalf("compiled DNS servers has type %T, want []map[string]any", compiled["servers"])
			}
			if len(servers) != 1 {
				t.Fatalf("compiled DNS servers length = %d, want 1", len(servers))
			}

			server := servers[0]
			if server["type"] != "https" {
				t.Fatalf("default DNS type = %v, want https", server["type"])
			}
			if server["tag"] != "default-dns" {
				t.Fatalf("default DNS tag = %v, want default-dns", server["tag"])
			}
			if server["server"] != "1.1.1.1" {
				t.Fatalf("default DNS server = %v, want 1.1.1.1", server["server"])
			}
			if compiled["final"] != "default-dns" {
				t.Fatalf("default DNS final = %v, want default-dns", compiled["final"])
			}
		})
	}
}

func TestCompileDNSRuleUsesUint32RewriteTTLAndTimeout(t *testing.T) {
	compiled, err := compileDNSRule(topology.DNSRule{
		Action:     "route",
		RewriteTTL: " 3600 ",
		Timeout:    "5s",
	})
	if err != nil {
		t.Fatal(err)
	}

	rewriteTTL, ok := compiled["rewrite_ttl"].(uint32)
	if !ok {
		t.Fatalf("rewrite_ttl type = %T, want uint32", compiled["rewrite_ttl"])
	}
	if rewriteTTL != 3600 {
		t.Fatalf("rewrite_ttl = %d, want 3600", rewriteTTL)
	}
	if timeout, exists := compiled["timeout"]; !exists || timeout != "5s" {
		t.Fatalf("timeout = %#v, want 5s", timeout)
	}
}

func TestCompileDNSRuleTimeoutActions(t *testing.T) {
	tests := []struct {
		action    string
		supported bool
	}{
		{action: "route", supported: true},
		{action: "evaluate", supported: true},
		{action: "route-options", supported: true},
		{action: "respond"},
		{action: "reject"},
		{action: "predefined"},
	}

	for _, test := range tests {
		t.Run(test.action, func(t *testing.T) {
			compiled, err := compileDNSRule(topology.DNSRule{
				Action:  test.action,
				Timeout: "5s",
			})
			if test.supported {
				if err != nil {
					t.Fatalf("compile DNS rule: %v", err)
				}
				if timeout, exists := compiled["timeout"]; !exists || timeout != "5s" {
					t.Fatalf("timeout = %#v, want 5s", timeout)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected unsupported timeout error, got rule %+v", compiled)
			}
			wantError := "dns rule action \"" + test.action + "\" does not support timeout"
			if err.Error() != wantError {
				t.Fatalf("compile DNS rule error = %q, want %q", err, wantError)
			}
		})
	}
}

func TestCompileDNSRuleRejectsInvalidRewriteTTL(t *testing.T) {
	for _, rewriteTTL := range []string{"-1", "1.5", "4294967296", "invalid"} {
		t.Run(rewriteTTL, func(t *testing.T) {
			_, err := compileDNSRule(topology.DNSRule{
				Action:     "route",
				RewriteTTL: rewriteTTL,
			})
			if err == nil {
				t.Fatal("expected invalid rewrite_ttl error")
			}
		})
	}
}
