# Embedded sing-box configuration layer

The Go sources in this directory are derived from sing-box and retain the
original package structure under an internal module path.

- Project: `github.com/SagerNet/sing-box`
- Source tag: `v1.14.0`
- Upstream commit: `0b8995879f29a9b98ee027bc17b75e101445b238`
- ACP fork commit: `fdd1af6da5f0763e10df5a5301d7e6150a4a9bb2`
- Imported packages: `schema`, `option`, `constant`, `constant/goos`,
  `common/badversion`, `experimental/locale`, and
  `experimental/deprecated`

The imported configuration-layer sources have no changes between the upstream
and ACP commits above. Imports referring to other sing-box packages were
mechanically rewritten to the corresponding internal paths in this module.

One behavior-neutral lint-only adjustment wraps the deferred
`buf.Put(buffer)` call in `option/dns_record.go` and explicitly discards its
error result.

The original license is preserved in [LICENSE](LICENSE).
