package inboundprovider

import (
	"encoding/json"
	"fmt"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
	sharedprovider "github.com/acp/node-agent/api/provider"
)

type BuildResult struct {
	Inbound  map[string]any
	Tag      string
	Protocol string
	Sniff    bool
}

type Adapter interface {
	ID() string
	Version() uint32
	Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error)
}

type Registry struct {
	adapters map[string]Adapter
}

func NewRegistry(adapters ...Adapter) *Registry {
	registry := &Registry{adapters: make(map[string]Adapter, len(adapters))}
	for _, adapter := range adapters {
		if adapter != nil {
			registry.adapters[adapter.ID()] = adapter
		}
	}
	return registry
}

func DefaultRegistry() *Registry {
	return NewRegistry(vlessRealityVisionAdapter{}, hysteria2SalamanderAdapter{})
}

func (r *Registry) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	if r == nil {
		return BuildResult{}, fmt.Errorf("inbound provider registry is not configured")
	}
	adapter, ok := r.adapters[node.ProviderID]
	if !ok {
		return BuildResult{}, fmt.Errorf("node %s has unsupported provider %q", node.NodeID, node.ProviderID)
	}
	if node.ProviderConfigVersion != adapter.Version() {
		return BuildResult{}, fmt.Errorf("node %s provider %s config version %d is unsupported", node.NodeID, node.ProviderID, node.ProviderConfigVersion)
	}
	return adapter.Build(node, users)
}

type vlessRealityVisionAdapter struct{}

func (vlessRealityVisionAdapter) ID() string      { return sharedprovider.VLESSRealityVisionID }
func (vlessRealityVisionAdapter) Version() uint32 { return sharedprovider.CurrentConfigVersion }

func (vlessRealityVisionAdapter) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	var cfg sharedprovider.VLESSRealityVisionConfig
	if err := json.Unmarshal(node.ProviderConfig, &cfg); err != nil {
		return BuildResult{}, fmt.Errorf("node %s decode vless provider config: %w", node.NodeID, err)
	}
	if cfg.ListenPort == 0 || cfg.TLS.Reality.PrivateKey == "" || len(cfg.TLS.Reality.ShortID) == 0 {
		return BuildResult{}, fmt.Errorf("node %s vless provider config is incomplete", node.NodeID)
	}
	tag := cfg.Tag
	if tag == "" {
		tag = node.NodeID
	}
	listen := cfg.Listen
	if listen == "" {
		listen = topology.DefaultInboundListen
	}
	compiledUsers := make([]map[string]any, 0, len(users))
	for _, user := range users {
		if user.Credential == "" {
			return BuildResult{}, fmt.Errorf("node %s user %s credential is required", node.NodeID, user.UserID)
		}
		entry := map[string]any{"name": userName(user), "uuid": user.Credential}
		if cfg.Flow != "" {
			entry["flow"] = cfg.Flow
		}
		compiledUsers = append(compiledUsers, entry)
	}
	inbound := map[string]any{
		"type": "vless", "tag": tag, "listen": listen, "listen_port": cfg.ListenPort,
		"users": compiledUsers,
		// VLESS inbounds always accept Mux connections. Clients independently
		// decide whether to use multiplexing on their outbounds.
		"multiplex": map[string]any{"enabled": true},
		"tls": map[string]any{
			"enabled":     cfg.TLS.Enabled,
			"server_name": cfg.TLS.ServerName,
			"reality": map[string]any{
				"enabled":     cfg.TLS.Reality.Enabled,
				"handshake":   map[string]any{"server": cfg.TLS.Reality.Handshake.Server, "server_port": cfg.TLS.Reality.Handshake.ServerPort},
				"private_key": cfg.TLS.Reality.PrivateKey,
				"short_id":    append([]string(nil), cfg.TLS.Reality.ShortID...),
			},
		},
	}
	if cfg.TCPFastOpen {
		inbound["tcp_fast_open"] = true
	}
	return BuildResult{Inbound: inbound, Tag: tag, Protocol: "vless", Sniff: cfg.Sniff}, nil
}

type hysteria2SalamanderAdapter struct{}

func (hysteria2SalamanderAdapter) ID() string      { return sharedprovider.Hysteria2SalamanderID }
func (hysteria2SalamanderAdapter) Version() uint32 { return sharedprovider.CurrentConfigVersion }

func (hysteria2SalamanderAdapter) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	var cfg sharedprovider.Hysteria2SalamanderConfig
	if err := json.Unmarshal(node.ProviderConfig, &cfg); err != nil {
		return BuildResult{}, fmt.Errorf("node %s decode hysteria2 provider config: %w", node.NodeID, err)
	}
	if cfg.ListenPort == 0 || cfg.TLS.CertificatePEM == "" || cfg.TLS.PrivateKeyPEM == "" {
		return BuildResult{}, fmt.Errorf("node %s hysteria2 provider config is incomplete", node.NodeID)
	}
	if cfg.UpMbps < 0 || cfg.DownMbps < 0 {
		return BuildResult{}, fmt.Errorf(
			"node %s hysteria2 bandwidth must not be negative: up_mbps=%d down_mbps=%d",
			node.NodeID, cfg.UpMbps, cfg.DownMbps,
		)
	}
	switch cfg.Obfs.Type {
	case "salamander":
		if cfg.Obfs.Password == "" {
			return BuildResult{}, fmt.Errorf("node %s hysteria2 salamander password is required", node.NodeID)
		}
		if cfg.Masquerade != nil {
			return BuildResult{}, fmt.Errorf("node %s hysteria2 masquerade requires obfs to be disabled", node.NodeID)
		}
	case "", "none":
		if cfg.Obfs.Password != "" {
			return BuildResult{}, fmt.Errorf("node %s hysteria2 obfs password requires salamander", node.NodeID)
		}
	default:
		return BuildResult{}, fmt.Errorf("node %s has unsupported hysteria2 obfs type %q", node.NodeID, cfg.Obfs.Type)
	}
	if cfg.Masquerade != nil {
		switch cfg.Masquerade.Type {
		case "proxy":
			if cfg.Masquerade.URL == "" {
				return BuildResult{}, fmt.Errorf("node %s hysteria2 proxy masquerade config is incomplete", node.NodeID)
			}
		case "string":
			if cfg.Masquerade.Content == "" {
				return BuildResult{}, fmt.Errorf("node %s hysteria2 fixed response masquerade content is required", node.NodeID)
			}
		default:
			return BuildResult{}, fmt.Errorf("node %s has unsupported hysteria2 masquerade type %q", node.NodeID, cfg.Masquerade.Type)
		}
	}
	tag := cfg.Tag
	if tag == "" {
		tag = node.NodeID
	}
	listen := cfg.Listen
	if listen == "" {
		listen = topology.DefaultInboundListen
	}
	compiledUsers := make([]map[string]any, 0, len(users))
	for _, user := range users {
		if user.Credential == "" {
			return BuildResult{}, fmt.Errorf("node %s user %s credential is required", node.NodeID, user.UserID)
		}
		compiledUsers = append(compiledUsers, map[string]any{"name": userName(user), "password": user.Credential})
	}
	inbound := map[string]any{
		"type": "hysteria2", "tag": tag, "listen": listen, "listen_port": cfg.ListenPort,
		"users": compiledUsers,
		"tls": map[string]any{
			"enabled":     true,
			"server_name": cfg.TLS.ServerName,
			"certificate": []string{cfg.TLS.CertificatePEM},
			"key":         []string{cfg.TLS.PrivateKeyPEM},
		},
	}
	if cfg.Obfs.Type == "salamander" {
		inbound["obfs"] = map[string]any{"type": "salamander", "password": cfg.Obfs.Password}
	}
	if cfg.Masquerade != nil {
		switch cfg.Masquerade.Type {
		case "proxy":
			inbound["masquerade"] = map[string]any{
				"type":         "proxy",
				"url":          cfg.Masquerade.URL,
				"rewrite_host": cfg.Masquerade.RewriteHost,
			}
		case "string":
			inbound["masquerade"] = map[string]any{
				"type":    "string",
				"headers": map[string]string{"Content-Type": "text/html; charset=utf-8"},
				"content": cfg.Masquerade.Content,
			}
		}
	}
	// sing-box maps these to SendBPS and ReceiveBPS independently and falls back
	// to BBR per direction, so a configured direction must never be dropped
	// because the other one is unset.
	if cfg.UpMbps > 0 {
		inbound["up_mbps"] = cfg.UpMbps
	}
	if cfg.DownMbps > 0 {
		inbound["down_mbps"] = cfg.DownMbps
	}
	return BuildResult{Inbound: inbound, Tag: tag, Protocol: "hysteria2"}, nil
}

func userName(user topology.UserCredential) string {
	if user.UserID != "" {
		return user.UserID
	}
	return user.Name
}
