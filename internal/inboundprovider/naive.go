package inboundprovider

import (
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
	"github.com/acp/node-agent/api/provider"
)

type naiveTLSAdapter struct{}

func (naiveTLSAdapter) ID() string      { return provider.NaiveTLSID }
func (naiveTLSAdapter) Version() uint32 { return provider.CurrentConfigVersion }
func (naiveTLSAdapter) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	cfg, err := provider.DecodeNaiveConfig(node.ProviderConfig)
	if err != nil {
		return BuildResult{}, err
	}
	if err = (option.WebOptions(cfg.Web)).Validate(); err != nil {
		return BuildResult{}, err
	}
	accounts := make([]provider.NaiveAccount, len(users))
	for i, user := range users {
		accounts[i] = provider.NaiveAccount{UserID: user.UserID, Credential: user.Credential}
	}
	tag, inbound, err := provider.CompileNaiveInbound(node.ProviderConfig, accounts)
	return BuildResult{Tag: tag, Protocol: "naive", Sniff: true, Inbound: inbound}, err
}
