# Compatibility

This module mirrors the panel-facing sing-box configuration behavior from
`github.com/acp/node-agent/pkg/singboxconfig` without importing the node-agent
runtime or sing-box protocol constructors.

Embedded configuration source provenance (unchanged):

- historical support release tag: `v1.14.0-acp.1`
- embedded sing-box configuration source tag: `v1.14.0`
- upstream sing-box source commit: `0b8995879f29a9b98ee027bc17b75e101445b238`
- original sing-box-acp source baseline: `e7ba3f961942ccf63e7173009498b16308eb93fc`
- original node-agent source baseline: `b77f8129bd48e67684498c26aaa2309e13f42ffb`
- sing-quic upstream commit: `4ab2eceaac81e073f53b22ec72ff56aaea89a3d1`
- sing-quic version: `v0.7.0-beta.4`

Current configuration compatibility target:

- node-agent commit: `d74fc89b40293b8f08cf066e822b5952a7f43c53`
- sing-box-acp commit: `b355c6e7280ab9c8d97eff071aa062ca1cb3f901`
- minimum Go version: `1.25.5`
- node-agent API module: `github.com/acp/node-agent`
- supported outbound types, in order: `direct`, `selector`, `urltest`,
  `shadowsocks`, `trojan`, `vless`, `hysteria2`

The embedded configuration-layer packages retain their original source
provenance, import-path adaptation, GPLv3 license, and source notice under
`internal/singbox`. The current fork's configuration-layer source was compared
against these packages and required no embedded-source update. The historical
`v1.14.0-acp.1` tag is not moved or redefined by this compatibility repair;
consumers must pin the new support commit to use the repaired compiler.

The panel-facing compiler now follows the current node's default Hysteria2
sniff rule (`300ms`), REALITY listener safety checks, DNS hostname bootstrap
selection and cycle avoidance, and default HTTP client selection for remote
rule sets. DNS rule validation also preserves the node's error precedence.
Compatibility is checked against the actual compiler behavior, not merely the
unchanged `pkg/singboxconfig` wrapper files. The parity suite covers DNS
bootstrap tag collisions and detours, default/explicit rule-set HTTP clients,
and invalid REALITY listener configurations in addition to the public APIs.

The root module has no direct `github.com/sagernet/sing-box` requirement or
replacement, and no production package or binary dependency on sing-box. The
isolated compatibility-test module owns the replacements needed to execute the
legacy node-agent implementation. An embedded source upgrade is complete only
after the four public API parity tests and the production dependency-boundary
check pass.

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
