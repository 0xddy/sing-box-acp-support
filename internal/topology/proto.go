package topology

import (
	"bytes"
	"encoding/json"
	"fmt"

	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FromSnapshot converts the public ACP wire model into the compiler's private
// topology model. Only snapshot compilation concerns live in this module.
func FromSnapshot(defaultMachineID string, snapshot *acpv1.TopologySnapshot) (MachineTopology, error) {
	top := MachineTopology{MachineID: defaultMachineID}
	if snapshot == nil {
		return top, nil
	}

	top.MachineID = snapshot.GetMachineId()
	top.Revision = snapshot.GetRevision()
	top.Nodes = make([]NodeInstance, 0, len(snapshot.GetNodes()))
	top.Outbounds = OutboundsFromProto(snapshot.GetOutbounds())
	var err error
	top.Route, err = RouteFromProto(snapshot.GetRoute())
	if err != nil {
		return MachineTopology{}, err
	}
	top.DNS, err = DNSFromProto(snapshot.GetDns())
	if err != nil {
		return MachineTopology{}, err
	}
	top.Snapshot = proto.Clone(snapshot).(*acpv1.TopologySnapshot)
	if top.MachineID == "" {
		top.MachineID = defaultMachineID
	}
	for _, node := range snapshot.GetNodes() {
		top.Nodes = append(top.Nodes, NodeFromProto(node))
	}
	return top, nil
}

func NodeFromProto(node *acpv1.NodeTopology) NodeInstance {
	if node == nil {
		return NodeInstance{}
	}
	users := make([]UserCredential, 0, len(node.GetUsers()))
	for _, user := range node.GetUsers() {
		users = append(users, UserFromProto(user))
	}
	return NodeInstance{
		NodeID:                node.GetNodeId(),
		ProviderID:            node.GetProviderId(),
		ProviderConfigVersion: node.GetProviderConfigVersion(),
		ProviderConfig:        append(json.RawMessage(nil), node.GetProviderConfigJson()...),
		Users:                 users,
	}
}

func UserFromProto(user *acpv1.UserCredential) UserCredential {
	if user == nil {
		return UserCredential{}
	}
	return UserCredential{
		UserID:                user.GetUserId(),
		Name:                  user.GetName(),
		Credential:            user.GetCredential(),
		Status:                userStatusFromProto(user.GetStatus()),
		UploadSpeedLimitBps:   user.GetUploadSpeedLimitBps(),
		DownloadSpeedLimitBps: user.GetDownloadSpeedLimitBps(),
	}
}

func OutboundsFromProto(outbounds []*acpv1.OutboundConfig) []Outbound {
	if len(outbounds) == 0 {
		return nil
	}
	converted := make([]Outbound, 0, len(outbounds))
	for _, outbound := range outbounds {
		if outbound == nil {
			continue
		}
		converted = append(converted, Outbound{
			Type:    outbound.GetType(),
			Tag:     outbound.GetTag(),
			Options: append(json.RawMessage(nil), outbound.GetOptionsJson()...),
		})
	}
	return converted
}

func RouteFromProto(route *acpv1.RouteConfig) (*Route, error) {
	if route == nil {
		return nil, nil
	}
	var converted Route
	if err := decodeProtoConfig(route, &converted); err != nil {
		return nil, fmt.Errorf("decode route configuration: %w", err)
	}
	return &converted, nil
}

func DNSFromProto(dns *acpv1.DNSConfig) (*DNS, error) {
	if dns == nil {
		return nil, nil
	}
	var converted DNS
	if err := decodeProtoConfig(dns, &converted); err != nil {
		return nil, fmt.Errorf("decode DNS configuration: %w", err)
	}
	return &converted, nil
}

func userStatusFromProto(status acpv1.UserStatus) string {
	if status == acpv1.UserStatus_USER_STATUS_DISABLED {
		return "disabled"
	}
	return ""
}

func decodeProtoConfig(src proto.Message, dst any) error {
	// protojson omits unknown wire fields. Reject them before conversion so a
	// removed or unsupported option cannot silently change routing behavior.
	if err := rejectUnknownProtoFields(src.ProtoReflect()); err != nil {
		return err
	}
	data, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(src)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func rejectUnknownProtoFields(message protoreflect.Message) error {
	if len(message.GetUnknown()) != 0 {
		return fmt.Errorf("%s contains unsupported protobuf fields", message.Descriptor().FullName())
	}
	var err error
	message.Range(func(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		switch {
		case field.IsMap() && field.MapValue().Message() != nil:
			value.Map().Range(func(_ protoreflect.MapKey, value protoreflect.Value) bool {
				err = rejectUnknownProtoFields(value.Message())
				return err == nil
			})
		case field.IsList() && field.Message() != nil:
			list := value.List()
			for i := 0; i < list.Len(); i++ {
				if err = rejectUnknownProtoFields(list.Get(i).Message()); err != nil {
					break
				}
			}
		case !field.IsMap() && field.Message() != nil:
			err = rejectUnknownProtoFields(value.Message())
		}
		return err == nil
	})
	return err
}
