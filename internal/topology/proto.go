package topology

import (
	"encoding/json"

	acpv1 "github.com/acp/node-agent/api/acp/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// FromSnapshot converts the public ACP wire model into the compiler's private
// topology model. Only snapshot compilation concerns live in this module.
func FromSnapshot(defaultMachineID string, snapshot *acpv1.TopologySnapshot) MachineTopology {
	top := MachineTopology{MachineID: defaultMachineID}
	if snapshot == nil {
		return top
	}

	top.MachineID = snapshot.GetMachineId()
	top.Revision = snapshot.GetRevision()
	top.Nodes = make([]NodeInstance, 0, len(snapshot.GetNodes()))
	top.Outbounds = OutboundsFromProto(snapshot.GetOutbounds())
	top.Route = RouteFromProto(snapshot.GetRoute())
	top.DNS = DNSFromProto(snapshot.GetDns())
	top.Snapshot = proto.Clone(snapshot).(*acpv1.TopologySnapshot)
	if top.MachineID == "" {
		top.MachineID = defaultMachineID
	}
	for _, node := range snapshot.GetNodes() {
		top.Nodes = append(top.Nodes, NodeFromProto(node))
	}
	return top
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

func RouteFromProto(route *acpv1.RouteConfig) *Route {
	if route == nil {
		return nil
	}
	var converted Route
	decodeProtoConfig(route, &converted)
	return &converted
}

func DNSFromProto(dns *acpv1.DNSConfig) *DNS {
	if dns == nil {
		return nil
	}
	var converted DNS
	decodeProtoConfig(dns, &converted)
	return &converted
}

func userStatusFromProto(status acpv1.UserStatus) string {
	if status == acpv1.UserStatus_USER_STATUS_DISABLED {
		return "disabled"
	}
	return ""
}

func decodeProtoConfig(src proto.Message, dst any) {
	data, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(src)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, dst)
}
