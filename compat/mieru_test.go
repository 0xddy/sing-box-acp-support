package compat

import (
	support "github.com/0xddy/sing-box-acp-support/singboxconfig"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"github.com/acp/node-agent/api/provider"
	"testing"
)

func TestMieruCompilationParity(t *testing.T) {
	node := &acpv1.NodeTopology{NodeId: "n", ProviderId: provider.MieruID, ProviderConfigVersion: 1, ProviderConfigJson: []byte(`{"type":"mieru","listen_port":2027,"transport":"udp","mtu":1400,"traffic_pattern":"CAc=","user_hint_is_mandatory":true}`)}
	snapshot := &acpv1.TopologySnapshot{MachineId: "m", Nodes: []*acpv1.NodeTopology{node}}
	for _, tc := range []struct {
		name    string
		users   []*acpv1.UserCredential
		invalid bool
	}{
		{"empty", nil, false},
		{"mixed", []*acpv1.UserCredential{{UserId: "42", Name: "display only", Credential: "account-secret", Status: acpv1.UserStatus_USER_STATUS_ACTIVE}, {UserId: "43", Status: acpv1.UserStatus_USER_STATUS_DISABLED}}, false},
		{"missing account credential", []*acpv1.UserCredential{{UserId: "34", Credential: "34", Status: acpv1.UserStatus_USER_STATUS_ACTIVE}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			node.Users = tc.users
			compareCompiledSnapshot(t, snapshot)
			if _, err := support.CompileSnapshot(snapshot); (err != nil) != tc.invalid {
				t.Fatalf("validation: %v", err)
			}
		})
	}
	node.Users = nil
	node.ProviderConfigJson = []byte(`{"type":"mieru","listen_port":2027,"traffic_pattern":"CA=="}`)
	compareCompiledSnapshot(t, snapshot)
	if _, err := support.CompileSnapshot(snapshot); err == nil {
		t.Fatal("support accepted an invalid Mieru profile")
	}
}
