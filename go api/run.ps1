# Скрипт запуска API сервера
# Использование: .\run.ps1

Write-Host "Запуск Go API сервера..." -ForegroundColor Green

# Опционально: можно установить CONFIG_PATH, если нужен другой конфиг
# $env:CONFIG_PATH = ".\config\prod.yaml"

# Запуск сервера
go run .\cmd\api\main.go
