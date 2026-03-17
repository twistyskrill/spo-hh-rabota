# Скрипт для сборки и запуска сервера
# Использование: .\start.ps1

Write-Host "Сборка и запуск сервера..." -ForegroundColor Cyan
Write-Host ""

# Сборка
.\build.ps1

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Запуск сервера..." -ForegroundColor Green
    Write-Host ""
    
    # Запуск собранного бинарника
    .\bin\api.exe
} else {
    Write-Host "Запуск отменён из-за ошибки сборки" -ForegroundColor Red
    exit 1
}
