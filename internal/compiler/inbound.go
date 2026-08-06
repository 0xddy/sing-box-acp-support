package compiler

import (
	"errors"

	"github.com/0xddy/sing-box-acp-support/internal/inboundprovider"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

var inboundProviders = inboundprovider.DefaultRegistry()

func compileInbound(node topology.NodeInstance) (inboundprovider.BuildResult, error) {
	if node.NodeID == "" {
		return inboundprovider.BuildResult{}, errors.New("node_id is required")
	}
	return inboundProviders.Build(node, activeUsers(node.Users))
}

func activeUsers(users []topology.UserCredential) []topology.UserCredential {
	out := make([]topology.UserCredential, 0, len(users))
	for _, user := range users {
		if user.Status == "disabled" {
			continue
		}
		out = append(out, user)
	}
	return out
}
