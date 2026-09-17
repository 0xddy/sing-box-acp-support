package topology

import (
	"encoding/json"

	acpv1 "github.com/acp/node-agent/api/acp/v1"
)

const (
	VLESSFlowRealityVision = "xtls-rprx-vision"
	DefaultInboundListen   = "::"
	OutboundTypeDirect     = "direct"
	DefaultDirectOutbound  = "direct"
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

type NodeInstance struct {
	NodeID                string           `json:"node_id"`
	ProviderID            string           `json:"provider_id"`
	ProviderConfigVersion uint32           `json:"provider_config_version"`
	ProviderConfig        json.RawMessage  `json:"provider_config"`
	Users                 []UserCredential `json:"users"`
}

type UserCredential struct {
	UserID                string `json:"user_id"`
	Name                  string `json:"name,omitempty"`
	Credential            string `json:"credential,omitempty"`
	Status                string `json:"status,omitempty"`
	UploadSpeedLimitBps   uint64 `json:"upload_speed_limit_bps,omitempty"`
	DownloadSpeedLimitBps uint64 `json:"download_speed_limit_bps,omitempty"`
}

type Outbound struct {
	Type    string          `json:"type"`
	Tag     string          `json:"tag"`
	Options json.RawMessage `json:"options,omitempty"`
}

type DirectActionOptions struct {
	BindInterface        string                `json:"bind_interface,omitempty"`
	Inet4BindAddress     string                `json:"inet4_bind_address,omitempty"`
	Inet6BindAddress     string                `json:"inet6_bind_address,omitempty"`
	RoutingMark          uint32                `json:"routing_mark,omitempty"`
	ReuseAddr            bool                  `json:"reuse_addr,omitempty"`
	ConnectTimeout       string                `json:"connect_timeout,omitempty"`
	TCPFastOpen          bool                  `json:"tcp_fast_open,omitempty"`
	TCPMultiPath         bool                  `json:"tcp_multi_path,omitempty"`
	UDPFragment          *bool                 `json:"udp_fragment,omitempty"`
	DomainStrategy       string                `json:"domain_strategy,omitempty"`
	BindAddressNoPort    bool                  `json:"bind_address_no_port,omitempty"`
	ProtectPath          string                `json:"protect_path,omitempty"`
	NetNS                string                `json:"netns,omitempty"`
	DisableTCPKeepAlive  bool                  `json:"disable_tcp_keep_alive,omitempty"`
	TCPKeepAlive         string                `json:"tcp_keep_alive,omitempty"`
	TCPKeepAliveInterval string                `json:"tcp_keep_alive_interval,omitempty"`
	DomainResolver       *DomainResolveOptions `json:"domain_resolver,omitempty"`
	NetworkStrategy      string                `json:"network_strategy,omitempty"`
	NetworkType          []string              `json:"network_type,omitempty"`
	FallbackNetworkType  []string              `json:"fallback_network_type,omitempty"`
	FallbackDelay        string                `json:"fallback_delay,omitempty"`
}

type DomainResolveOptions struct {
	Server       string  `json:"server,omitempty"`
	Strategy     string  `json:"strategy,omitempty"`
	DisableCache bool    `json:"disable_cache,omitempty"`
	RewriteTTL   *uint32 `json:"rewrite_ttl,omitempty"`
	ClientSubnet string  `json:"client_subnet,omitempty"`
}

type Route struct {
	Rules                      []RouteRule           `json:"rules,omitempty"`
	RuleSets                   []RouteRuleSet        `json:"rule_sets,omitempty"`
	Final                      string                `json:"final,omitempty"`
	AutoDetectInterface        bool                  `json:"auto_detect_interface,omitempty"`
	DefaultInterface           string                `json:"default_interface,omitempty"`
	DefaultMark                uint32                `json:"default_mark,omitempty"`
	FindProcess                bool                  `json:"find_process,omitempty"`
	GeoIP                      *GeoIPOptions         `json:"geoip,omitempty"`
	Geosite                    *GeositeOptions       `json:"geosite,omitempty"`
	OverrideAndroidVPN         bool                  `json:"override_android_vpn,omitempty"`
	DefaultDomainResolver      *DomainResolveOptions `json:"default_domain_resolver,omitempty"`
	DefaultNetworkStrategy     string                `json:"default_network_strategy,omitempty"`
	DefaultNetworkType         []string              `json:"default_network_type,omitempty"`
	DefaultFallbackNetworkType []string              `json:"default_fallback_network_type,omitempty"`
	DefaultFallbackDelay       string                `json:"default_fallback_delay,omitempty"`
}

type DNS struct {
	Rules   []DNSRule   `json:"rules,omitempty"`
	Servers []DNSServer `json:"servers,omitempty"`
	Final   string      `json:"final,omitempty"`
}

type DNSServer struct {
	Type   string `json:"type,omitempty"`
	Tag    string `json:"tag,omitempty"`
	Server string `json:"server,omitempty"`
	Detour string `json:"detour,omitempty"`
}

type DNSRule struct {
	Inbound       []string `json:"inbound,omitempty"`
	Domain        []string `json:"domain,omitempty"`
	DomainSuffix  []string `json:"domain_suffix,omitempty"`
	DomainKeyword []string `json:"domain_keyword,omitempty"`
	DomainRegex   []string `json:"domain_regex,omitempty"`
	RuleSet       []string `json:"rule_set,omitempty"`
	Action        string   `json:"action,omitempty"`
	RCode         string   `json:"rcode,omitempty"`
	Server        string   `json:"server,omitempty"`
	Method        string   `json:"method,omitempty"`
	NoDrop        bool     `json:"no_drop,omitempty"`
	Answer        []string `json:"answer,omitempty"`
	NS            []string `json:"ns,omitempty"`
	Extra         []string `json:"extra,omitempty"`
	DisableCache  bool     `json:"disable_cache,omitempty"`
	RewriteTTL    string   `json:"rewrite_ttl,omitempty"`
	Timeout       string   `json:"timeout,omitempty"`
	ClientSubnet  string   `json:"client_subnet,omitempty"`
}

type GeoIPOptions struct {
	Path           string `json:"path,omitempty"`
	DownloadURL    string `json:"download_url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty"`
}

type GeositeOptions struct {
	Path           string `json:"path,omitempty"`
	DownloadURL    string `json:"download_url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty"`
}

type RouteRule struct {
	Type                     string                `json:"type,omitempty"`
	Inbound                  []string              `json:"inbound,omitempty"`
	Network                  []string              `json:"network,omitempty"`
	IPVersion                uint8                 `json:"ip_version,omitempty"`
	Domain                   []string              `json:"domain,omitempty"`
	DomainSuffix             []string              `json:"domain_suffix,omitempty"`
	DomainKeyword            []string              `json:"domain_keyword,omitempty"`
	DomainRegex              []string              `json:"domain_regex,omitempty"`
	SourceIPCidr             []string              `json:"source_ip_cidr,omitempty"`
	IPCidr                   []string              `json:"ip_cidr,omitempty"`
	SourceIPIsPrivate        *bool                 `json:"source_ip_is_private,omitempty"`
	IPIsPrivate              *bool                 `json:"ip_is_private,omitempty"`
	Port                     []uint32              `json:"port,omitempty"`
	PortRange                []string              `json:"port_range,omitempty"`
	SourcePort               []uint32              `json:"source_port,omitempty"`
	SourcePortRange          []string              `json:"source_port_range,omitempty"`
	Protocol                 []string              `json:"protocol,omitempty"`
	RuleSet                  []string              `json:"rule_set,omitempty"`
	Invert                   bool                  `json:"invert,omitempty"`
	Action                   string                `json:"action,omitempty"`
	Outbound                 string                `json:"outbound,omitempty"`
	Method                   string                `json:"method,omitempty"`
	NoDrop                   bool                  `json:"no_drop,omitempty"`
	Mode                     string                `json:"mode,omitempty"`
	Rules                    []RouteRule           `json:"rules,omitempty"`
	AuthUser                 []string              `json:"auth_user,omitempty"`
	Client                   []string              `json:"client,omitempty"`
	Geosite                  []string              `json:"geosite,omitempty"`
	SourceGeoIP              []string              `json:"source_geoip,omitempty"`
	GeoIP                    []string              `json:"geoip,omitempty"`
	ProcessName              []string              `json:"process_name,omitempty"`
	ProcessPath              []string              `json:"process_path,omitempty"`
	ProcessPathRegex         []string              `json:"process_path_regex,omitempty"`
	PackageName              []string              `json:"package_name,omitempty"`
	User                     []string              `json:"user,omitempty"`
	UserID                   []int32               `json:"user_id,omitempty"`
	ClashMode                string                `json:"clash_mode,omitempty"`
	NetworkType              []string              `json:"network_type,omitempty"`
	NetworkIsExpensive       *bool                 `json:"network_is_expensive,omitempty"`
	NetworkIsConstrained     *bool                 `json:"network_is_constrained,omitempty"`
	WIFISSID                 []string              `json:"wifi_ssid,omitempty"`
	WIFIBSSID                []string              `json:"wifi_bssid,omitempty"`
	DefaultInterfaceAddress  []string              `json:"default_interface_address,omitempty"`
	PreferredBy              []string              `json:"preferred_by,omitempty"`
	RuleSetIPCIDRMatchSource bool                  `json:"rule_set_ip_cidr_match_source,omitempty"`
	RouteOptions             *RouteActionOptions   `json:"route_options,omitempty"`
	DirectOptions            *DirectActionOptions  `json:"direct_options,omitempty"`
	SniffOptions             *SniffActionOptions   `json:"sniff_options,omitempty"`
	ResolveOptions           *ResolveActionOptions `json:"resolve_options,omitempty"`
}

type RouteRuleSet struct {
	Type           string         `json:"type"`
	Tag            string         `json:"tag"`
	Format         string         `json:"format,omitempty"`
	Path           string         `json:"path,omitempty"`
	URL            string         `json:"url,omitempty"`
	DownloadDetour string         `json:"download_detour,omitempty"`
	UpdateInterval string         `json:"update_interval,omitempty"`
	Rules          []HeadlessRule `json:"rules,omitempty"`
}

type RouteActionOptions struct {
	OverrideAddress           string `json:"override_address,omitempty"`
	OverridePort              uint32 `json:"override_port,omitempty"`
	NetworkStrategy           string `json:"network_strategy,omitempty"`
	FallbackDelay             uint32 `json:"fallback_delay,omitempty"`
	UDPDisableDomainUnmapping bool   `json:"udp_disable_domain_unmapping,omitempty"`
	UDPConnect                bool   `json:"udp_connect,omitempty"`
	UDPTimeout                string `json:"udp_timeout,omitempty"`
	TLSFragment               bool   `json:"tls_fragment,omitempty"`
	TLSFragmentFallbackDelay  string `json:"tls_fragment_fallback_delay,omitempty"`
	TLSRecordFragment         bool   `json:"tls_record_fragment,omitempty"`
}

type SniffActionOptions struct {
	Sniffer []string `json:"sniffer,omitempty"`
	Timeout string   `json:"timeout,omitempty"`
}

type ResolveActionOptions struct {
	Server       string  `json:"server,omitempty"`
	Strategy     string  `json:"strategy,omitempty"`
	DisableCache bool    `json:"disable_cache,omitempty"`
	RewriteTTL   *uint32 `json:"rewrite_ttl,omitempty"`
	ClientSubnet string  `json:"client_subnet,omitempty"`
}

type HeadlessRule struct {
	Type                    string         `json:"type,omitempty"`
	Network                 []string       `json:"network,omitempty"`
	Domain                  []string       `json:"domain,omitempty"`
	DomainSuffix            []string       `json:"domain_suffix,omitempty"`
	DomainKeyword           []string       `json:"domain_keyword,omitempty"`
	DomainRegex             []string       `json:"domain_regex,omitempty"`
	SourceIPCidr            []string       `json:"source_ip_cidr,omitempty"`
	IPCidr                  []string       `json:"ip_cidr,omitempty"`
	SourcePort              []uint32       `json:"source_port,omitempty"`
	SourcePortRange         []string       `json:"source_port_range,omitempty"`
	Port                    []uint32       `json:"port,omitempty"`
	PortRange               []string       `json:"port_range,omitempty"`
	ProcessName             []string       `json:"process_name,omitempty"`
	ProcessPath             []string       `json:"process_path,omitempty"`
	ProcessPathRegex        []string       `json:"process_path_regex,omitempty"`
	PackageName             []string       `json:"package_name,omitempty"`
	NetworkType             []string       `json:"network_type,omitempty"`
	NetworkIsExpensive      *bool          `json:"network_is_expensive,omitempty"`
	NetworkIsConstrained    *bool          `json:"network_is_constrained,omitempty"`
	WIFISSID                []string       `json:"wifi_ssid,omitempty"`
	WIFIBSSID               []string       `json:"wifi_bssid,omitempty"`
	DefaultInterfaceAddress []string       `json:"default_interface_address,omitempty"`
	Invert                  bool           `json:"invert,omitempty"`
	Mode                    string         `json:"mode,omitempty"`
	Rules                   []HeadlessRule `json:"rules,omitempty"`
}
