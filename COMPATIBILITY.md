# Compatibility

This module mirrors the panel-facing sing-box configuration behavior from
`github.com/acp/node-agent/pkg/singboxconfig` without importing the node-agent
runtime or sing-box protocol constructors.

Compatibility baseline:

- planned support release tag: `v1.14.0-acp.1`
- embedded sing-box configuration source tag: `v1.14.0`
- upstream sing-box source commit: `0b8995879f29a9b98ee027bc17b75e101445b238`
- sing-box-acp commit: `e7ba3f961942ccf63e7173009498b16308eb93fc`
- node-agent commit: `b77f8129bd48e67684498c26aaa2309e13f42ffb`
- sing-quic upstream commit: `4ab2eceaac81e073f53b22ec72ff56aaea89a3d1`
- sing-quic version: `v0.7.0-beta.4`
- minimum Go version: `1.25.5`
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
