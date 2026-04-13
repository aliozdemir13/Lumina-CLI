$BinaryName = "lumina.exe"

switch ($args[0]) {
    "build" {
        go build -o $BinaryName main.go
        Write-Host "✅ Build complete: $BinaryName" -ForegroundColor Cyan
    }
    "test" {
        go test -v ./internal/
    }
    "cover" {
        go test -cover ./internal/
    }
    "clean" {
        if (Test-Path $BinaryName) {
            Remove-Item $BinaryName
            Write-Host "🗑️ Artifacts removed." -ForegroundColor Yellow
        }
    }
    "run" {
        go build -o $BinaryName main.go
        .\$BinaryName --help
    }
    Default {
        Write-Host "Usage: .\build.ps1 {build|test|cover|clean|run}" -ForegroundColor White
    }
}