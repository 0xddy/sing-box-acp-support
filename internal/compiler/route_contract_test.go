package compiler

import (
	"context"
	"encoding/json"
	"testing"

	C "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	singjson "github.com/sagernet/sing/common/json"
	"google.golang.org/protobuf/proto"
)

func TestRouteContractMatchesEmbeddedSingBox(t *testing.T) {
	route := &acpv1.RouteConfig{
		DefaultNetworkStrategy: "fallback", DefaultNetworkType: []string{"ethernet"},
		DefaultFallbackNetworkType: []string{"wifi"}, DefaultFallbackDelay: "300ms",
		Rules: []*acpv1.RouteRule{
			{Action: "route-options", RouteOptions: &acpv1.RouteActionOptions{NetworkStrategy: "hybrid", FallbackDelay: 300}},
			{Action: "direct", DirectOptions: &acpv1.DirectActionOptions{
				NetworkStrategy: "default", NetworkType: []string{"wifi"},
				FallbackNetworkType: []string{"cellular"}, FallbackDelay: "500ms", UdpFragment: proto.Bool(false),
			}},
		},
	}
	top, err := topology.FromSnapshot("machine", &acpv1.TopologySnapshot{Route: route})
	if err != nil {
		t.Fatal(err)
	}
	compiled, err := Compile(top)
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		Route json.RawMessage `json:"route"`
	}
	if err = json.Unmarshal(compiled, &config); err != nil {
		t.Fatal(err)
	}
	parsed, err := singjson.UnmarshalExtendedContext[option.RouteOptions](context.Background(), config.Route)
	if err != nil {
		t.Fatalf("embedded sing-box rejected compiled route: %v\n%s", err, config.Route)
	}
	if parsed.DefaultNetworkStrategy == nil || C.NetworkStrategy(*parsed.DefaultNetworkStrategy) != C.NetworkStrategyFallback {
		t.Fatalf("default strategy was lost: %+v", parsed.DefaultNetworkStrategy)
	}
	direct := parsed.Rules[1].DefaultOptions.DirectOptions
	if direct.UDPFragment == nil || *direct.UDPFragment {
		t.Fatalf("explicit udp_fragment=false was lost: %+v", direct.UDPFragment)
	}
}

func TestCompileRejectsUnsupportedNetworkStrategy(t *testing.T) {
	for _, route := range []*topology.Route{
		{DefaultNetworkStrategy: "invalid"},
		{Rules: []topology.RouteRule{{Action: "route-options", RouteOptions: &topology.RouteActionOptions{NetworkStrategy: "invalid"}}}},
		{Rules: []topology.RouteRule{{Action: "direct", DirectOptions: &topology.DirectActionOptions{NetworkStrategy: "invalid"}}}},
	} {
		if _, err := Compile(topology.MachineTopology{MachineID: "machine", Route: route}); err == nil {
			t.Fatal("unsupported network strategy was accepted")
		}
	}
}
