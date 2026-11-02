package inboundprovider

import (
	"fmt"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
	"github.com/acp/node-agent/api/provider"
)

type juicityAdapter struct{}

func (juicityAdapter) ID() string      { return provider.JuicityID }
func (juicityAdapter) Version() uint32 { return provider.CurrentConfigVersion }

func (juicityAdapter) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	accounts := make([]provider.JuicityAccount, len(users))
	for i, user := range users {
		accounts[i] = provider.JuicityAccount{UserID: user.UserID, Credential: user.Credential}
	}
	tag, inbound, err := provider.CompileJuicityInbound(node.ProviderConfig, accounts)
	if err != nil {
		return BuildResult{}, fmt.Errorf("node %s: %w", node.NodeID, err)
	}
	return BuildResult{Tag: tag, Protocol: "juicity", Inbound: inbound, Sniff: true}, nil
}
