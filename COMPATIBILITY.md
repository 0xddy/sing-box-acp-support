# Compatibility

This module mirrors the panel-facing sing-box configuration behavior from
`github.com/acp/node-agent/pkg/singboxconfig` without importing the node-agent
runtime or sing-box protocol constructors.

Compatibility baseline:

- support release tag: `v1.13.19-acp.2`
- embedded sing-box configuration source tag: `v1.13.19`
- upstream sing-box source commit: `b5ebaa1fc0f2b94256180b95468e73ef53caa27d`
- sing-box-acp commit: `bc3b80fcf35f2d41f743c42805898a25b322c73e`
- node-agent commit: `b7907ed95ea63c048761063bde2023a7386aee33`
- sing-quic upstream commit: `d83826c306d7c008cafb3fd6d8ee07cbcfd656ed`
- sing-quic pseudo-version: `v0.6.4-0.20260803041914-d83826c306d7`
- node-agent API module: `github.com/acp/node-agent`
- supported outbound types, in order: `direct`, `selector`, `urltest`,
  `shadowsocks`, `trojan`, `vless`, `hysteria2`

The embedded configuration-layer packages are identical between the upstream
source commit and the ACP fork commit listed above. Their original GPLv3
license and source notice are preserved under `internal/singbox`.

The root module has no direct `github.com/sagernet/sing-box` requirement or
replacement, and no production package or binary dependency on sing-box. The
isolated compatibility-test module owns the replacements needed to execute the
legacy node-agent implementation. An embedded source upgrade is complete only
after the four public API parity tests and the production dependency-boundary
check pass.

`go list -m all` can still display upstream sing-box module metadata inherited
from the `github.com/acp/node-agent` contract module. That module remains the
source of the public `acpv1.TopologySnapshot` type, but none of its sing-box
packages are compiled into this support module. Removing that metadata would
require splitting or duplicating the ACP contract module and would change the
public Go type identity.

The support package performs the same option-decoding validation as the
current panel integration. It does not instantiate protocols and therefore
does not promise that a configuration can start successfully in a particular
runtime environment.
