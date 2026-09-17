# Embedded sing-box configuration layer

The Go sources in this directory are derived from sing-box and retain the
original package structure under an internal module path.

- Project: `github.com/SagerNet/sing-box`
- Source tag: `v1.14.1`
- Upstream commit: `1ac1a339cb1223e9c70eae14c44411c75033c02d`
- ACP fork commit: `302319a74e006d426d460a2ebc13c79e314d5c1f`
- Imported packages: `schema`, `option`, `constant`, `constant/goos`,
  `common/badversion`, `experimental/locale`, and
  `experimental/deprecated`

The imported configuration layer includes the upstream `v1.14.1` changes to
Clash API, HTTP/2, and OpenVPN option types. The ACP fork adds one configuration
change in `option/rule_nested.go`: explicit actions are decoded by their own
option types so a direct action's duration-valued `fallback_delay` is accepted.
Imports referring to other sing-box packages were mechanically rewritten to
the corresponding internal paths in this module.

One behavior-neutral lint-only adjustment wraps the deferred
`buf.Put(buffer)` call in `option/dns_record.go` and explicitly discards its
error result.

The original license is preserved in [LICENSE](LICENSE).
