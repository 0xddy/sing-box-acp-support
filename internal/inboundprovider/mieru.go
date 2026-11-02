package inboundprovider

import (
	"fmt"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
	"github.com/acp/node-agent/api/provider"
)

type mieruAdapter struct{}

func (mieruAdapter) ID() string      { return provider.MieruID }
func (mieruAdapter) Version() uint32 { return provider.CurrentConfigVersion }
func (mieruAdapter) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	accounts := make([]provider.MieruUser, len(users))
	for i, user := range users {
		accounts[i] = provider.MieruUser{Name: user.UserID, Password: user.Credential}
	}
	tag, inbound, err := provider.CompileMieruInbound(node.ProviderConfig, node.NodeID, accounts)
	if err != nil {
		return BuildResult{}, fmt.Errorf("node %s: %w", node.NodeID, err)
	}
	return BuildResult{Tag: tag, Protocol: "mieru", Inbound: inbound, Sniff: true}, nil
}
