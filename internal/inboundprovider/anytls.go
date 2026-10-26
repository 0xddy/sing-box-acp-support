package inboundprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
	"github.com/0xddy/sing-box-acp-support/internal/topology"
	"github.com/acp/node-agent/api/provider"
)

type anytlsTLSAdapter struct{}

func (anytlsTLSAdapter) ID() string      { return provider.AnyTLSTLSID }
func (anytlsTLSAdapter) Version() uint32 { return provider.CurrentConfigVersion }
func (anytlsTLSAdapter) Build(node topology.NodeInstance, users []topology.UserCredential) (BuildResult, error) {
	var cfg provider.AnyTLSTLSConfig
	decoder := json.NewDecoder(bytes.NewReader(node.ProviderConfig))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return BuildResult{}, fmt.Errorf("node %s decode AnyTLS config: %w", node.NodeID, err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return BuildResult{}, fmt.Errorf("node %s AnyTLS config must contain exactly one JSON object", node.NodeID)
	}
	if cfg.Type != "anytls" || cfg.Tag == "" || cfg.Listen == "" || cfg.ListenPort == 0 || !cfg.TLS.Enabled || cfg.TLS.ServerName == "" || cfg.TLS.CertificatePEM == "" || cfg.TLS.PrivateKeyPEM == "" {
		return BuildResult{}, fmt.Errorf("node %s AnyTLS config is incomplete", node.NodeID)
	}
	if err := (option.WebOptions(cfg.Web)).Validate(); err != nil {
		return BuildResult{}, fmt.Errorf("node %s web: %w", node.NodeID, err)
	}
	compiledUsers := make([]map[string]any, 0, len(users))
	for _, user := range users {
		if user.UserID == "" || user.Credential == "" {
			return BuildResult{}, fmt.Errorf("node %s AnyTLS user ID and credential are required", node.NodeID)
		}
		compiledUsers = append(compiledUsers, map[string]any{"name": user.UserID, "password": user.Credential})
	}
	return BuildResult{Tag: cfg.Tag, Protocol: "anytls", Sniff: true, Inbound: map[string]any{
		"type": "anytls", "tag": cfg.Tag, "listen": cfg.Listen, "listen_port": cfg.ListenPort, "users": compiledUsers, "web": cfg.Web,
		"tls": map[string]any{"enabled": true, "server_name": cfg.TLS.ServerName, "certificate": []string{cfg.TLS.CertificatePEM}, "key": []string{cfg.TLS.PrivateKeyPEM}, "alpn": []string{"h2", "http/1.1"}},
	}}, nil
}
