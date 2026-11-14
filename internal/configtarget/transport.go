package configtarget

import (
	"encoding/json"

	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
)

func (Schema) DecodeTransport(raw json.RawMessage) (json.RawMessage, error) {
	var transport option.V2RayTransportOptions
	if err := json.Unmarshal(raw, &transport); err != nil {
		return nil, err
	}
	return json.Marshal(transport)
}
