param(
    [string]$SourceRoot = ""
)

$ErrorActionPreference = "Stop"

$projectRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot ".."))
$destinationRoot = [System.IO.Path]::GetFullPath((Join-Path $projectRoot "internal\singbox"))
$expectedDestination = [System.IO.Path]::GetFullPath((Join-Path $projectRoot "internal\singbox"))
if ($destinationRoot -ne $expectedDestination) {
    throw "refusing to sync outside the expected destination: $destinationRoot"
}

if ([string]::IsNullOrWhiteSpace($SourceRoot)) {
    $SourceRoot = Join-Path $projectRoot "..\third_party\sing-box-acp"
}
$sourcePath = (Resolve-Path -LiteralPath $SourceRoot).Path
if (-not (Test-Path -LiteralPath (Join-Path $sourcePath "option\options.go"))) {
    throw "source does not look like a sing-box tree: $sourcePath"
}
if (-not (Test-Path -LiteralPath (Join-Path $sourcePath "LICENSE"))) {
    throw "source license is missing: $sourcePath"
}

$packages = @(
    "option",
    "constant",
    "constant\goos",
    "common\badversion",
    "experimental\deprecated",
    "experimental\locale"
)

$rewrites = [ordered]@{
    "github.com/sagernet/sing-box/constant/goos" = "github.com/0xddy/sing-box-acp-support/internal/singbox/constant/goos"
    "github.com/sagernet/sing-box/common/badversion" = "github.com/0xddy/sing-box-acp-support/internal/singbox/common/badversion"
    "github.com/sagernet/sing-box/experimental/deprecated" = "github.com/0xddy/sing-box-acp-support/internal/singbox/experimental/deprecated"
    "github.com/sagernet/sing-box/experimental/locale" = "github.com/0xddy/sing-box-acp-support/internal/singbox/experimental/locale"
    "github.com/sagernet/sing-box/constant" = "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
    "github.com/sagernet/sing-box/option" = "github.com/0xddy/sing-box-acp-support/internal/singbox/option"
}

$utf8WithoutBom = New-Object System.Text.UTF8Encoding($false)
foreach ($package in $packages) {
    $sourcePackage = Join-Path $sourcePath $package
    if (-not (Test-Path -LiteralPath $sourcePackage)) {
        throw "source package is missing: $sourcePackage"
    }

    $destinationPackage = Join-Path $destinationRoot $package
    New-Item -ItemType Directory -Path $destinationPackage -Force | Out-Null
    Get-ChildItem -LiteralPath $destinationPackage -File -Filter "*.go" |
        Remove-Item -Force

    Get-ChildItem -LiteralPath $sourcePackage -File -Filter "*.go" |
        Where-Object { $_.Name -notlike "*_test.go" } |
        ForEach-Object {
            $destinationFile = Join-Path $destinationPackage $_.Name
            $content = [System.IO.File]::ReadAllText($_.FullName)
            foreach ($rewrite in $rewrites.GetEnumerator()) {
                $content = $content.Replace($rewrite.Key, $rewrite.Value)
            }

            # Local lint-only adjustment documented in internal/singbox/NOTICE.md.
            if ($package -eq "option" -and $_.Name -eq "dns_record.go") {
                $content = $content.Replace(
                    "defer buf.Put(buffer)",
                    "defer func() { _ = buf.Put(buffer) }()"
                )
            }
            [System.IO.File]::WriteAllText($destinationFile, $content, $utf8WithoutBom)
        }
}

Copy-Item -LiteralPath (Join-Path $sourcePath "LICENSE") -Destination (Join-Path $destinationRoot "LICENSE") -Force

$externalImports = Get-ChildItem -LiteralPath $destinationRoot -Recurse -File -Filter "*.go" |
    Select-String -Pattern "github.com/sagernet/sing-box/"
if ($externalImports) {
    $externalImports | ForEach-Object { Write-Error $_.Line }
    throw "sync left external sing-box imports in the embedded source"
}

& gofmt -w $destinationRoot
if ($LASTEXITCODE -ne 0) {
    throw "gofmt failed"
}

$sourceCommit = (& git -C $sourcePath rev-parse HEAD 2>$null)
Write-Host "Synchronized sing-box configuration packages from $sourcePath"
if ($LASTEXITCODE -eq 0 -and $sourceCommit) {
    Write-Host "Source commit: $sourceCommit"
}
Write-Host "Review and update internal/singbox/NOTICE.md and COMPATIBILITY.md when the source baseline changes."
