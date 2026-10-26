# Embedded sing-box configuration layer

The Go sources in this directory are derived from sing-box and retain the
original package structure under an internal module path.

- Project: `github.com/SagerNet/sing-box`
- Source tag: `v1.14.1`
- Upstream commit: `1ac1a339cb1223e9c70eae14c44411c75033c02d`
- ACP fork commit: `2e3c5327aa358ebf3592ba1953f299b762672030`
- Imported packages: `schema`, `option`, `constant`, `constant/goos`,
  `common/badversion`, `experimental/locale`, and
  `experimental/deprecated`

The imported configuration layer includes the upstream `v1.14.1` changes to
Clash API, HTTP/2, and OpenVPN option types. The ACP fork adds one configuration
change in `option/rule_nested.go`: explicit actions are decoded by their own
option types so a direct action's duration-valued `fallback_delay` is accepted.
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
