# Embedded sing-box configuration layer

The Go sources in this directory are derived from sing-box and retain the
original package structure under an internal module path.

- Project: `github.com/SagerNet/sing-box`
- Source tag: `v1.13.19`
- Upstream commit: `b5ebaa1fc0f2b94256180b95468e73ef53caa27d`
- ACP fork commit: `bc3b80fcf35f2d41f743c42805898a25b322c73e`
- Imported packages: `option`, `constant`, `constant/goos`,
  `common/badversion`, `experimental/locale`, and
  `experimental/deprecated`

The imported configuration-layer sources have no changes between the upstream
and ACP commits above. Imports referring to other sing-box packages were
mechanically rewritten to the corresponding internal paths in this module.

One behavior-neutral lint-only adjustment wraps the deferred
`buf.Put(buffer)` call in `option/dns_record.go` and explicitly discards its
error result.

The original license is preserved in [LICENSE](LICENSE).
