# Compatibility

This module mirrors the panel-facing sing-box configuration behavior from
`github.com/acp/node-agent/pkg/singboxconfig` without importing the node-agent
runtime or sing-box protocol constructors.

Compatibility baselines:

- historical support release tag: `v1.14.0-acp.1`
- upgrade branch: `agent/upgrade-sing-box-1.14.1`
- planned support release tag: `v1.14.1-acp.1`
- embedded sing-box configuration source tag: `v1.14.1`
- upstream sing-box source commit: `1ac1a339cb1223e9c70eae14c44411c75033c02d`
- node-agent commit: `e058ad2dfab6d9b3881c6e93e8225723d6a02851`
- sing-box-acp commit: `4c07687ad9180f8c6b83fc5e4fd5a73a7ea7b8d7`
- sing-quic upstream commit: `4ab2eceaac81e073f53b22ec72ff56aaea89a3d1`
- sing-quic version: `v0.7.0` (the same source commit as `v0.7.0-beta.4`)

Current configuration compatibility target:

- minimum Go version: `1.25.5`
- node-agent API module: `github.com/acp/node-agent`
- supported outbound types, in order: `direct`, `selector`, `urltest`,
  `shadowsocks`, `trojan`, `vless`, `hysteria2`

The embedded configuration-layer packages retain their original source
provenance, import-path adaptation, GPLv3 license, and source notice under
`internal/singbox`. This upgrade includes the three upstream option-model
changes present in the embedded package set and the ACP fork's direct-action
decoding fix. The historical `v1.14.0-acp.1` tag and
`agent/upgrade-sing-box-1.14.0` branch remain available for rollback.

The panel-facing compiler now follows the current node's default Hysteria2
sniff rule (`300ms`), REALITY listener safety checks, DNS hostname bootstrap
selection and cycle avoidance, and default HTTP client selection for remote
rule sets. DNS rule validation also preserves the node's error precedence.
The route contract uses the sing-box string-valued network strategy enum,
rejects stale protobuf fields, and preserves explicit `udp_fragment=false`.
REALITY validation rejects malformed SNI, private keys, and short IDs before
the runtime receives them.
Compatibility is checked against the actual compiler behavior through the
public `pkg/singboxconfig` API. The parity suite covers DNS
bootstrap tag collisions and detours, default/explicit rule-set HTTP clients,
invalid REALITY listener configurations, and route option decoding in addition
to the public APIs.

The root module has no direct `github.com/sagernet/sing-box` requirement or
replacement, and no production package or binary dependency on sing-box. The
isolated compatibility-test module owns the replacements needed to execute the
node-agent implementation. An embedded source upgrade is complete only after
the compatibility parity suite and the production dependency-boundary check
pass.

Validation commands (using the pinned sibling node-agent and sing-box-acp
checkouts, with no network listeners or database connections):

```sh
go test -mod=readonly ./...
cd compat
go test -mod=readonly -tags with_utls ./...
```

The root suite includes the production dependency-boundary check. No new
module dependency or runtime constructor is needed by this repair.

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
