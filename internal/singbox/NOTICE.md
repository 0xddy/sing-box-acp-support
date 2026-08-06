# Embedded sing-box configuration layer

The Go sources in this directory are derived from sing-box and retain the
original package structure under an internal module path.

- Project: `github.com/SagerNet/sing-box`
- Source tag: `v1.13.16`
- Upstream commit: `17ec3c71af8ca946dc50bf0d927c39fc77322aec`
- ACP fork commit: `bf3f82a697c888bdaf41b927c3e72728e4e19909`
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
