# Embedded sing-box configuration layer

The Go sources in this directory are derived from sing-box and retain the
original package structure under an internal module path.

- Project: `github.com/SagerNet/sing-box`
- Source tag: `v1.14.4`
- Upstream commit: `c005f35300ce76427dc285481ec5695f5b179cd2`
- ACP fork commit: `373f2e593d5a28326a98a6038356482d9ef0f6b9`
- Imported packages: `schema`, `option`, `constant`, `constant/goos`,
  `common/badversion`, `experimental/locale`, and
  `experimental/deprecated`

The imported configuration layer tracks the v1.14.4 option types.
Imports referring to other sing-box packages were mechanically rewritten to
the corresponding internal paths in this module.

The ACP AnyTLS integration adds `option/web.go` and the `web` field in
`option/anytls.go`. These configuration-only sources are synchronized from the
sibling ACP core using `tools/sync_singbox_config.ps1`; the HTTP runtime is not
copied into this module. The commit recorded above identifies the synchronized
configuration source.

One behavior-neutral lint-only adjustment wraps the deferred
`buf.Put(buffer)` call in `option/dns_record.go` and explicitly discards its
error result.

The original license is preserved in [LICENSE](LICENSE).

The XHTTP transport options in `option/v2ray_xhttp.go` and its transport
discriminator are mirrored from ACP core `c01553ca9`. The adapted option
types are covered by the accompanying MPL-2.0 notice and license.

Juicity's option types and protocol discriminator are mirrored from ACP core
`59aa816bf56670f5658c8b2bdc6299bab1eed3f4`.
