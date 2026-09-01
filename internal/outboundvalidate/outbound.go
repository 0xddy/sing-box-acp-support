// Package outboundvalidate provides sing-box option decoding without importing
// protocol constructors or runtime registries.
package outboundvalidate

import (
	"bytes"
	"context"
	"fmt"
	"reflect"

	C "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
	singjson "github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/service"
)

type registration struct {
	outboundType string
	create       func() any
}

var registrations = []registration{
	{outboundType: C.TypeDirect, create: func() any { return new(option.DirectOutboundOptions) }},
	{outboundType: C.TypeSelector, create: func() any { return new(option.SelectorOutboundOptions) }},
	{outboundType: C.TypeURLTest, create: func() any { return new(option.URLTestOutboundOptions) }},
	{outboundType: C.TypeShadowsocks, create: func() any { return new(option.ShadowsocksOutboundOptions) }},
	{outboundType: C.TypeTrojan, create: func() any { return new(option.TrojanOutboundOptions) }},
	{outboundType: C.TypeVLESS, create: func() any { return new(option.VLESSOutboundOptions) }},
	{outboundType: C.TypeHysteria2, create: func() any { return new(option.Hysteria2OutboundOptions) }},
}

type optionsRegistry struct {
	factories map[string]func() any
}

func newOptionsRegistry() *optionsRegistry {
	factories := make(map[string]func() any, len(registrations))
	for _, item := range registrations {
		factories[item.outboundType] = item.create
	}
	return &optionsRegistry{factories: factories}
}

func (r *optionsRegistry) CreateOptions(outboundType string) (any, bool) {
	create, ok := r.factories[outboundType]
	if !ok {
		return nil, false
	}
	return create(), true
}

func (r *optionsRegistry) OptionTypes() []string {
	types := make([]string, 0, len(registrations))
	for _, item := range registrations {
		types = append(types, item.outboundType)
	}
	return types
}

var registry = newOptionsRegistry()

func validationContext() context.Context {
	return service.ContextWith[option.OutboundOptionsRegistry](context.Background(), registry)
}

// Validate decodes a complete sing-box outbound object with the option types
// selected by this support module.
func Validate(encoded []byte) error {
	_, err := singjson.UnmarshalExtendedContext[option.Outbound](validationContext(), encoded)
	return err
}

// SupportedTypes returns a copy of the ordered supported type list.
func SupportedTypes() []string {
	types := make([]string, 0, len(registrations))
	for _, item := range registrations {
		types = append(types, item.outboundType)
	}
	return types
}

// DirectOutboundIsEmpty mirrors the empty-dialer calculation used by the
// sing-box Direct outbound constructor.
func DirectOutboundIsEmpty(optionsJSON []byte) (bool, error) {
	rawOptions := bytes.TrimSpace(optionsJSON)
	if len(rawOptions) == 0 {
		rawOptions = []byte("{}")
	}
	options, err := singjson.UnmarshalExtendedContext[option.DirectOutboundOptions](validationContext(), rawOptions)
	if err != nil {
		return false, fmt.Errorf("invalid Direct outbound options: %w", err)
	}
	options.AbstractDialerOptions.UDPFragmentDefault = true
	return reflect.DeepEqual(
		options.DialerOptions,
		option.DialerOptions{
			AbstractDialerOptions: option.AbstractDialerOptions{UDPFragmentDefault: true},
		},
	), nil
}
