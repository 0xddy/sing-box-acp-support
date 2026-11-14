package topology

import (
	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"github.com/acp/node-agent/pkg/configcompile/model"
)

const (
	VLESSFlowRealityVision = model.VLESSFlowRealityVision
	DefaultInboundListen   = model.DefaultInboundListen
	OutboundTypeDirect     = model.OutboundTypeDirect
	DefaultDirectOutbound  = model.DefaultDirectOutbound
)

type MachineTopology struct {
	MachineID string                  `json:"machine_id"`
	Revision  uint64                  `json:"revision"`
	Nodes     []NodeInstance          `json:"nodes"`
	Outbounds []Outbound              `json:"outbounds,omitempty"`
	Route     *Route                  `json:"route,omitempty"`
	DNS       *DNS                    `json:"dns,omitempty"`
	Snapshot  *acpv1.TopologySnapshot `json:"-"`
}

type NodeInstance = model.NodeInstance
type UserCredential = model.UserCredential
type Outbound = model.Outbound
type DirectActionOptions = model.DirectActionOptions
type DomainResolveOptions = model.DomainResolveOptions
type Route = model.Route
type DNS = model.DNS
type DNSServer = model.DNSServer
type DNSRule = model.DNSRule
type GeoIPOptions = model.GeoIPOptions
type GeositeOptions = model.GeositeOptions
type RouteRule = model.RouteRule
type RouteRuleSet = model.RouteRuleSet
type RouteActionOptions = model.RouteActionOptions
type SniffActionOptions = model.SniffActionOptions
type ResolveActionOptions = model.ResolveActionOptions
type HeadlessRule = model.HeadlessRule

func (t MachineTopology) CompilerInput() model.MachineTopology {
	return model.MachineTopology{MachineID: t.MachineID, Revision: t.Revision, Nodes: t.Nodes, Outbounds: t.Outbounds, Route: t.Route, DNS: t.DNS}
}
