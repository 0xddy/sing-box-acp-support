// Package singboxconfig compiles ACP topology snapshots and validates the
// subset of sing-box outbound types supported by sing-box-acp.
package singboxconfig

import (
	"errors"

	"github.com/0xddy/sing-box-acp-support/internal/compiler"
	"github.com/0xddy/sing-box-acp-support/internal/outboundvalidate"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
	acpv1 "github.com/acp/node-agent/api/acp/v1"
)

// CompileSnapshot compiles the complete sing-box JSON represented by an ACP
// topology snapshot.
func CompileSnapshot(snapshot *acpv1.TopologySnapshot) ([]byte, error) {
	if snapshot == nil {
		return nil, errors.New("topology snapshot is required")
	}
	top, err := topology.FromSnapshot(snapshot.GetMachineId(), snapshot)
	if err != nil {
		return nil, err
	}
	return compiler.Compile(top)
}

// ValidateOutbound checks an outbound's raw options against the sing-box
// option types pinned by this support module.
func ValidateOutbound(outbound *acpv1.OutboundConfig) error {
	if outbound == nil {
		return errors.New("outbound is required")
	}
	return compiler.ValidateOutbound(topology.Outbound{
		Type:    outbound.GetType(),
		Tag:     outbound.GetTag(),
		Options: append([]byte(nil), outbound.GetOptionsJson()...),
	})
}

// SupportedOutboundTypes returns a copy of the ordered outbound type list.
func SupportedOutboundTypes() []string {
	return outboundvalidate.SupportedTypes()
}

// DirectOutboundIsEmpty reports whether sing-box treats a Direct outbound as
// an empty dialer, which cannot safely be used as a detour target.
func DirectOutboundIsEmpty(optionsJSON []byte) (bool, error) {
	return outboundvalidate.DirectOutboundIsEmpty(optionsJSON)
}
