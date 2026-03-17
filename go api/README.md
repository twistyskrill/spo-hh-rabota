# Go API - Сервис для поиска мастеров и заказчиков

API сервис для платформы поиска мастеров и размещения объявлений о работе.

## 🚀 Быстрый старт

### Запуск проекта

Теперь для запуска **не нужно** устанавливать переменную окружения `CONFIG_PATH`!

**Вариант 1: Запуск через Go (для разработки)**
```powershell
.\run.ps1
# или
go run .\cmd\api\main.go
```

**Вариант 2: Сборка и запуск бинарника**
```powershell
.\start.ps1
# или вручную:
.\build.ps1
.\bin\api.exe
```

**Вариант 3: Только сборка**
```powershell
.\build.ps1
```

### Конфигурация

По умолчанию используется файл `./config/local.yaml`

Если нужен другой конфиг, установите переменную окружения:
```powershell
$env:CONFIG_PATH = ".\config\prod.yaml"
go run .\cmd\api\main.go
```

## 📋 Требования

- **Go** 1.21+
- **PostgreSQL** 14+
- **Git**

## 🗄️ Настройка базы данных

1. Создайте базу данных PostgreSQL:
```sql
CREATE DATABASE handyman;
```

2. Настройте параметры подключения в `config/local.yaml`:
```yaml
db:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "ваш_пароль"
  dbname: "handyman"
  sslmode: "disable"
```

3. Миграции выполнятся автоматически при первом запуске (GORM AutoMigrate)

## 📚 Документация API

Полная документация API находится в файле [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

**Основные эндпоинты:**

### Аутентификация
- `POST /auth/register` - Регистрация
- `POST /auth/login` - Вход
- `GET /profile` - Получить профиль
- `PATCH /profile` - Обновить профиль

### Объявления клиентов
- `GET /ads` - Список объявлений (публичный)
- `GET /ads/{id}` - Получить объявление по ID
- `GET /my-ads` - Мои объявления
- `POST /my-ads` - Создать объявление
- `PATCH /my-ads/{id}` - Обновить объявление
- `DELETE /my-ads/{id}` - Удалить объявление

### Отклики мастеров
- `GET /responses` - Мои отклики (для мастеров)
- `POST /responses` - Откликнуться на объявление
- `DELETE /responses/{id}` - Отменить отклик

### Админ-панель
- `GET /admin/users` - Управление пользователями
- `GET /admin/ads` - Модерация объявлений
- `GET /admin/responses` - Модерация откликов
- `GET /admin/stats` - Статистика платформы
- `GET /admin/blacklist` - Черный список

> 📖 Полная документация админ-панели: [ADMIN_GUIDE.md](ADMIN_GUIDE.md)

### Мастера и справочники
- `GET /handyman` - Список мастеров
- `GET /handyman/{id}` - Мастер по ID
- `GET /info/categories` - Список категорий
- `GET /info/price-units` - Единицы измерения цены

## 🛠️ Структура проекта

```
.
├── cmd/
│   └── api/
│       └── main.go           # Точка входа
├── config/
│   └── local.yaml            # Конфигурация
├── internal/
│   ├── auth/                 # JWT и хеширование паролей
│   ├── config/               # Загрузка конфигурации
│   ├── handlers/             # HTTP handlers
│   │   ├── admin/           # Админ-панель
│   │   ├── ads/             # Объявления клиентов
│   │   ├── auth/            # Аутентификация и профиль
│   │   ├── info/            # Справочная информация
│   │   ├── sys/             # Системные эндпоинты
│   │   └── worker/          # Мастера
│   ├── middleware/          # Middleware (auth)
│   ├── models/              # Модели данных (GORM)
│   └── storage/             # Слой работы с БД
├── bin/                     # Скомпилированные бинарники
├── run.ps1                  # Скрипт запуска (dev)
├── build.ps1               # Скрипт сборки
├── start.ps1               # Сборка + запуск
└── API_DOCUMENTATION.md    # API документация

```

## 🔧 Разработка

### Установка зависимостей
```powershell
go mod download
```

### Проверка кода
```powershell
go vet ./...
go fmt ./...
```

### Сборка для production
```powershell
go build -ldflags="-s -w" -o bin/api.exe ./cmd/api
```

## 🌍 Переменные окружения

| Переменная | Описание | По умолчанию |
|-----------|----------|-------------|
| `CONFIG_PATH` | Путь к конфигу | `./config/local.yaml` |
| `ENV` | Окружение (local/dev/prod) | Из конфига |

## 📝 Примеры использования

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "Иван Иванов",
    "password": "secure123",
    "role": 1
  }'
```

### Вход
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure123"
  }'
```

### Создание объявления
```bash
curl -X POST http://localhost:8080/my-ads \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Требуется сантехник",
    "price": 3000,
    "category_id": 1,
    "price_unit_id": 1,
    "location": "Москва",
    "schedule": "будние дни, 9-18"
  }'
```

### Просмотр объявлений
```bash
# Все объявления (публичный доступ)
curl -X GET "http://localhost:8080/ads?category=1&location=Москва&limit=10"

# Мои объявления (требуется авторизация)
curl -X GET http://localhost:8080/my-ads \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Отклик мастера на объявление
```bash
# Создать отклик
curl -X POST http://localhost:8080/responses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer MASTER_TOKEN" \
  -d '{
    "ad_id": 5,
    "message": "Готов выполнить работу качественно и в срок",
    "proposed_price": 2500
  }'

# Мои отклики
curl -X GET http://localhost:8080/responses \
  -H "Authorization: Bearer MASTER_TOKEN"

# Отменить отклик
curl -X DELETE http://localhost:8080/responses/123 \
  -H "Authorization: Bearer MASTER_TOKEN"
```

## 🐛 Отладка

### Логи
Логи выводятся в stdout. Уровень логирования зависит от `env` в конфиге:
- `local` - DEBUG (текстовый формат)
- `dev` - DEBUG (JSON формат)
- `prod` - INFO (JSON формат)

### Частые проблемы

**Ошибка подключения к БД:**
- Проверьте, что PostgreSQL запущен
- Проверьте параметры в `config/local.yaml`
- Убедитесь, что база данных создана

**Порт уже занят:**
- Измените `address` в `config/local.yaml`
- Или завершите процесс на порту 8080

## 📄 Лицензия

Учебный проект

## 👨‍💻 Автор

СПО проект
