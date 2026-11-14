package inboundprovider

import (
	"github.com/0xddy/sing-box-acp-support/internal/configtarget"
	"github.com/acp/node-agent/pkg/configcompile/inbound"
)

type BuildResult = inbound.BuildResult
type Adapter = inbound.Adapter
type Registry = inbound.Registry

func NewRegistry(adapters ...Adapter) *Registry { return inbound.NewRegistry(adapters...) }
func DefaultRegistry() *Registry                { return inbound.DefaultRegistry(configtarget.Schema{}) }
