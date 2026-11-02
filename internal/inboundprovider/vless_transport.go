package inboundprovider

import (
	"encoding/json"
	"fmt"

	C "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
)

func vlessTransport(raw json.RawMessage) (*option.V2RayTransportOptions, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var header struct {
		Type             string           `json:"type"`
		DownloadSettings *json.RawMessage `json:"download_settings"`
		Xmux             *json.RawMessage `json:"xmux"`
	}
	if err := json.Unmarshal(raw, &header); err != nil {
		return nil, err
	}
	switch header.Type {
	case "", "raw":
		return nil, nil
	case C.V2RayTransportTypeXHTTP:
	default:
		return nil, fmt.Errorf("unsupported type %q (expected raw or xhttp)", header.Type)
	}
	if header.DownloadSettings != nil || header.Xmux != nil {
		return nil, fmt.Errorf("download_settings and xmux are XHTTP outbound options")
	}
	var transport option.V2RayTransportOptions
	if err := json.Unmarshal(raw, &transport); err != nil {
		return nil, err
	}
	opts := transport.XHTTPOptions
	switch opts.Mode {
	case "", "auto", "packet-up", "stream-up", "stream-one":
	default:
		return nil, fmt.Errorf("unsupported XHTTP mode %q", opts.Mode)
	}
	return &transport, nil
}
