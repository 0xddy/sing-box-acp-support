package configtarget

import "github.com/0xddy/sing-box-acp-support/internal/outboundvalidate"

func (Schema) ValidateOutbound(encoded []byte) error { return outboundvalidate.Validate(encoded) }
func (Schema) DirectOutboundIsEmpty(raw []byte) (bool, error) {
	return outboundvalidate.DirectOutboundIsEmpty(raw)
}
