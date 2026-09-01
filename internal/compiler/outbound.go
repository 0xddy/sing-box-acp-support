package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/0xddy/sing-box-acp-support/internal/outboundvalidate"
	C "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func compileOutbounds(outbounds []topology.Outbound) ([]map[string]any, map[string]struct{}, string, error) {
	if len(outbounds) == 0 {
		outbounds := []map[string]any{
			{
				"type": topology.OutboundTypeDirect,
				"tag":  topology.DefaultDirectOutbound,
			},
		}
		return outbounds, map[string]struct{}{topology.DefaultDirectOutbound: {}}, topology.DefaultDirectOutbound, nil
	}

	compiled := make([]map[string]any, 0, len(outbounds))
	tags := make(map[string]struct{}, len(outbounds))
	defaultRouteFinal := ""
	for _, outbound := range outbounds {
		if outbound.Type == "" {
			return nil, nil, "", errors.New("outbound type is required")
		}
		if outbound.Tag == "" {
			return nil, nil, "", fmt.Errorf("outbound %s tag is required", outbound.Type)
		}
		if _, ok := tags[outbound.Tag]; ok {
			return nil, nil, "", fmt.Errorf("duplicate outbound tag %q", outbound.Tag)
		}
		entry, err := compileOutbound(outbound)
		if err != nil {
			return nil, nil, "", err
		}
		compiled = append(compiled, entry)
		tags[outbound.Tag] = struct{}{}
		if outbound.Tag == topology.DefaultDirectOutbound {
			defaultRouteFinal = topology.DefaultDirectOutbound
		}
	}
	return compiled, tags, defaultRouteFinal, nil
}

// ValidateOutbound validates one topology outbound with the same path used by
// full machine configuration compilation.
func ValidateOutbound(outbound topology.Outbound) error {
	_, _, _, err := compileOutbounds([]topology.Outbound{outbound})
	return err
}

func compileOutbound(outbound topology.Outbound) (map[string]any, error) {
	rawOptions := bytes.TrimSpace(outbound.Options)
	if len(rawOptions) == 0 {
		rawOptions = []byte("{}")
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(rawOptions, &fields); err != nil {
		return nil, fmt.Errorf("outbound %q options must be a JSON object: %w", outbound.Tag, err)
	}
	if fields == nil {
		return nil, fmt.Errorf("outbound %q options must be a JSON object", outbound.Tag)
	}
	if outbound.Type == C.TypeHysteria2 && !hasJSONField(fields, "disable_chrome_parrot") {
		fields["disable_chrome_parrot"] = json.RawMessage("true")
	}

	entry := make(map[string]any, len(fields)+2)
	entry["type"] = outbound.Type
	entry["tag"] = outbound.Tag
	for key, value := range fields {
		if strings.EqualFold(key, "type") || strings.EqualFold(key, "tag") {
			return nil, fmt.Errorf("outbound %q options must not contain managed field %q", outbound.Tag, key)
		}
		entry[key] = append(json.RawMessage(nil), value...)
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("encode outbound %q: %w", outbound.Tag, err)
	}
	if err := outboundvalidate.Validate(encoded); err != nil {
		return nil, fmt.Errorf("invalid %s outbound %q options: %w", outbound.Type, outbound.Tag, err)
	}
	return entry, nil
}

func hasJSONField(fields map[string]json.RawMessage, name string) bool {
	for field := range fields {
		if strings.EqualFold(field, name) {
			return true
		}
	}
	return false
}
