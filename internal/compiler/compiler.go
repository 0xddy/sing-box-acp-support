package compiler

import (
	"github.com/0xddy/sing-box-acp-support/internal/configtarget"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
	"github.com/acp/node-agent/pkg/configcompile"
)

var engine = configcompile.New(configtarget.Schema{})

func Compile(top topology.MachineTopology) ([]byte, error) {
	return engine.Compile(top.CompilerInput())
}
func CompileDetailed(top topology.MachineTopology) (configcompile.Result, error) {
	return engine.CompileDetailed(top.CompilerInput())
}

func ValidateOutbound(outbound topology.Outbound) error { return engine.ValidateOutbound(outbound) }
func DirectOutboundIsEmpty(raw []byte) (bool, error) {
	return (configtarget.Schema{}).DirectOutboundIsEmpty(raw)
}
func compileOutbound(outbound topology.Outbound) (map[string]any, error) {
	return engine.CompileOutbound(outbound)
}
func compileDNS(dns *topology.DNS) (map[string]any, error) { return engine.CompileDNS(dns) }
func compileDNSRule(rule topology.DNSRule) (map[string]any, error) {
	return engine.CompileDNSRule(rule)
}
func compileRouteRuleSet(rule topology.RouteRuleSet) (map[string]any, error) {
	return configcompile.CompileRouteRuleSet(rule)
}
