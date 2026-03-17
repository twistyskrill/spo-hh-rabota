# Скрипт сборки проекта
# Использование: .\build.ps1

Write-Host "🔨 Сборка проекта..." -ForegroundColor Cyan

# Создаём директорию bin если её нет
if (-not (Test-Path ".\bin")) {
    New-Item -ItemType Directory -Path ".\bin" | Out-Null
}

# Сборка
go build -o .\bin\api.exe .\cmd\api\main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Сборка успешна! Исполняемый файл: .\bin\api.exe" -ForegroundColor Green
    
    # Показываем размер файла
    $fileSize = (Get-Item ".\bin\api.exe").Length / 1MB
    Write-Host "📦 Размер: $($fileSize.ToString('0.00')) MB" -ForegroundColor Gray
} else {
    Write-Host "❌ Ошибка сборки!" -ForegroundColor Red
    exit 1
}
