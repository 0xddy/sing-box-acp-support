package topology

import (
	"strings"
	"testing"

	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRouteConversionRejectsUnrepresentableFields(t *testing.T) {
	route := &acpv1.RouteConfig{Rules: []*acpv1.RouteRule{{Action: "reject", IpVersion: 300}}}
	if converted, err := RouteFromProto(route); err == nil || converted != nil {
		t.Fatalf("invalid ip_version must not become an unrestricted reject: converted=%+v error=%v", converted, err)
	}
	if _, err := FromSnapshot("machine", &acpv1.TopologySnapshot{Route: route}); err == nil {
		t.Fatal("snapshot conversion swallowed the route error")
	}
}

func TestRouteConversionRejectsUnsupportedDirectOptions(t *testing.T) {
	for _, fieldNumber := range []protowire.Number{1, 11} {
		direct := &acpv1.DirectActionOptions{}
		wire := protowire.AppendTag(nil, fieldNumber, protowire.BytesType)
		direct.ProtoReflect().SetUnknown(protowire.AppendString(wire, "unsupported"))
		converted, err := RouteFromProto(&acpv1.RouteConfig{Rules: []*acpv1.RouteRule{{Action: "direct", DirectOptions: direct}}})
		if err == nil || converted != nil || !strings.Contains(err.Error(), "unsupported protobuf fields") {
			t.Fatalf("field %d: converted=%+v error=%v", fieldNumber, converted, err)
		}
	}
}

func TestProtoConversionRejectsJSONModelDrift(t *testing.T) {
	source, err := structpb.NewStruct(map[string]any{"unexpected_field": true})
	if err != nil {
		t.Fatal(err)
	}
	if err = decodeProtoConfig(source, &Route{}); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("JSON model mismatch must fail: %v", err)
	}
}

func TestDNSConversionPropagatesMarshalErrors(t *testing.T) {
	dns := &acpv1.DNSConfig{Final: string([]byte{0xff})}
	if converted, err := DNSFromProto(dns); err == nil || converted != nil {
		t.Fatalf("invalid UTF-8 was accepted: converted=%+v error=%v", converted, err)
	}
	if _, err := FromSnapshot("machine", &acpv1.TopologySnapshot{Dns: dns}); err == nil {
		t.Fatal("snapshot conversion swallowed the DNS error")
	}
}
