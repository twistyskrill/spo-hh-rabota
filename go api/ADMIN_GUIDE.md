# Руководство по админ-панели

Администраторская панель для управления платформой поиска мастеров.

## 🔐 Требования доступа

Для доступа к админ-панели необходимо:
1. **JWT токен** (получается через `/auth/login`)
2. **Роль администратора** (`admin` или `administrator` в таблице `roles`)

## 📋 Доступные эндпоинты

Все эндпоинты требуют заголовок:
```
Authorization: Bearer YOUR_ADMIN_TOKEN
```

### Управление пользователями

#### Получить список всех пользователей
```http
GET /admin/users?limit=10&offset=0&role=user&search=ivan
```

**Query параметры:**
- `limit` - количество записей (по умолчанию: 10)
- `offset` - смещение для пагинации
- `role` - фильтр по роли (опционально)
- `search` - поиск по email или имени (опционально)

**Ответ:**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван Иванов",
      "phone": "+79001234567",
      "role_id": 1,
      "role_name": "user",
      "created_at": "2026-03-01T10:00:00Z",
      "have_worker_profile": true,
      "ads_count": 5,
      "responses_count": 12
    }
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

#### Получить пользователя по ID
```http
GET /admin/users/123
```

**Ответ:**
```json
{
  "id": 123,
  "email": "user@example.com",
  "name": "Иван Иванов",
  "phone": "+79001234567",
  "role": "user",
  "role_id": 1,
  "created_at": "2026-03-01T10:00:00Z",
  "worker_profile": {...},
  "ads_count": 5,
  "responses_count": 12
}
```

#### Удалить пользователя (soft delete)
```http
DELETE /admin/users/123
```

**Ответ:**
```json
{
  "message": "user deleted successfully"
}
```

#### Изменить роль пользователя
```http
PATCH /admin/users/123/role
Content-Type: application/json

{
  "role_name": "admin"
}
```

**Ответ:**
```json
{
  "message": "role updated successfully",
  "user_id": 123,
  "role": "admin"
}
```

---

### Модерация объявлений

> Статусы: `pending` (na рассмотрении) → `approved` (одобрено) / `rejected` (отклонено)

#### Получить все объявления
```http
GET /admin/ads?limit=10&offset=0&status=pending&category=Сантехника&user_id=5
```

**Query параметры:**
- `limit` - количество записей
- `offset` - смещение
- `status` - фильтр по статусу (`pending`, `approved`, `rejected`)
- `category` - фильтр по названию категории
- `user_id` - фильтр по автору объявления

**Ответ:**
```json
{
  "ads": [
    {
      "id": 15,
      "title": "Требуется сантехник",
      "price": 3000,
      "location": "Москва",
      "created_at": "2026-03-08T12:00:00Z",
      "category_name": "Сантехника",
      "price_unit_name": "час",
      "user_id": 5,
      "user_name": "Петр Петров",
      "user_email": "petr@example.com",
      "responses_count": 3,
      "status": "pending"
    }
  ],
  "total": 50,
  "limit": 10,
  "offset": 0
}
```

#### Одобрить объявление
```http
PATCH /admin/ads/15/approve
```

**Ответ:**
```json
{
  "message": "ad approved successfully",
  "ad_id": 15,
  "status": "approved"
}
```

#### Отклонить объявление
```http
PATCH /admin/ads/15/reject
```

**Ответ:**
```json
{
  "message": "ad rejected successfully",
  "ad_id": 15,
  "status": "rejected"
}
```

#### Удалить объявление
```http
DELETE /admin/ads/15
```

**Ответ:**
```json
{
  "message": "ad deleted successfully"
}
```

---

### Модерация профилей мастеров

#### Получить профили на модерации
```http
GET /admin/workers?status=pending&limit=10&offset=0
```

**Query параметры:**
- `status` - фильтр по статусу (`pending`, `approved`, `rejected`); по умолчанию `pending`
- `limit` / `offset` - пагинация

**Ответ:**
```json
{
  "workers": [
    {
      "user_id": 12,
      "name": "Сергей Мастеров",
      "email": "master@example.com",
      "phone": "+79001234567",
      "exp_years": 5,
      "description": "Профессиональный сантехник",
      "location": "Москва",
      "schedule": "Пн-Пт 9:00-18:00",
      "status": "pending"
    }
  ],
  "total": 8,
  "limit": 10,
  "offset": 0,
  "status": "pending"
}
```

#### Одобрить профиль мастера
```http
PATCH /admin/workers/12/approve
```

**Ответ:**
```json
{
  "message": "worker profile approved successfully",
  "worker_id": 12,
  "status": "approved"
}
```

#### Отклонить профиль мастера
```http
PATCH /admin/workers/12/reject
```

**Ответ:**
```json
{
  "message": "worker profile rejected successfully",
  "worker_id": 12,
  "status": "rejected"
}
```

---

### Модерация откликов

#### Получить все отклики
```http
GET /admin/responses?limit=10&offset=0&status=pending&worker_id=7
```

**Query параметры:**
- `limit` - количество записей
- `offset` - смещение
- `status` - фильтр по статусу (pending, accepted, rejected, cancelled)
- `worker_id` - фильтр по ID мастера

**Ответ:**
```json
{
  "responses": [
    {
      "id": 42,
      "ad_id": 15,
      "ad_title": "Требуется сантехник",
      "worker_id": 7,
      "worker_name": "Сергей Мастеров",
      "worker_email": "master@example.com",
      "message": "Готов выполнить работу",
      "proposed_price": 2500,
      "status": "pending",
      "created_at": "2026-03-08T14:00:00Z"
    }
  ],
  "total": 30,
  "limit": 10,
  "offset": 0
}
```

#### Удалить отклик
```http
DELETE /admin/responses/42
```

**Ответ:**
```json
{
  "message": "response deleted successfully"
}
```

---

### Статистика

#### Получить общую статистику платформы
```http
GET /admin/stats
```

**Ответ:**
```json
{
  "total_users": 150,
  "total_ads": 75,
  "total_responses": 230,
  "total_workers": 45,
  "users_by_role": {
    "user": 120,
    "admin": 5,
    "master": 25
  },
  "ads_by_category": {
    "Сантехника": 20,
    "Электрика": 15,
    "Ремонт": 40
  },
  "recent_activity": {
    "new_users_today": 5,
    "new_ads_today": 8,
    "new_responses_today": 15
  }
}
```

---

### Черный список

#### Получить черный список email
```http
GET /admin/blacklist
```

**Ответ:**
```json
{
  "blacklist": [
    {
      "email": "spam@example.com"
    },
    {
      "email": "blocked@test.com"
    }
  ],
  "total": 2
}
```

#### Добавить email в черный список
```http
POST /admin/blacklist
Content-Type: application/json

{
  "email": "spam@example.com"
}
```

**Ответ:**
```json
{
  "message": "email added to blacklist",
  "email": "spam@example.com"
}
```

#### Удалить email из черного списка
```http
DELETE /admin/blacklist/spam@example.com
```

**Ответ:**
```json
{
  "message": "email removed from blacklist",
  "email": "spam@example.com"
}
```

---

### Управление справочниками

#### Создать категорию
```http
POST /admin/categories
Content-Type: application/json

{
  "name": "Ландшафтный дизайн"
}
```

**Ответ:**
```json
{
  "message": "Category created successfully",
  "id": 10,
  "name": "Ландшафтный дизайн"
}
```

#### Обновить категорию
```http
PATCH /admin/categories/10
Content-Type: application/json

{
  "name": "Ландшафтный дизайн и озеленение"
}
```

**Ответ:**
```json
{
  "message": "Category updated successfully",
  "id": 10,
  "name": "Ландшафтный дизайн и озеленение"
}
```

#### Удалить категорию
```http
DELETE /admin/categories/10
```

**Ответ:**
```json
{
  "message": "Category deleted successfully"
}
```

**Примечание:** Категорию нельзя удалить, если она используется в объявлениях или профилях мастеров (возвращает код 409).

#### Создать единицу цены
```http
POST /admin/price-units
Content-Type: application/json

{
  "name": "за м²"
}
```

**Ответ:**
```json
{
  "message": "Price unit created successfully",
  "id": 5,
  "name": "за м²"
}
```

#### Обновить единицу цены
```http
PATCH /admin/price-units/5
Content-Type: application/json

{
  "name": "за квадратный метр"
}
```

**Ответ:**
```json
{
  "message": "Price unit updated successfully",
  "id": 5,
  "name": "за квадратный метр"
}
```

#### Удалить единицу цены
```http
DELETE /admin/price-units/5
```

**Ответ:**
```json
{
  "message": "Price unit deleted successfully"
}
```

**Примечание:** Единицу цены нельзя удалить, если она используется в объявлениях (возвращает код 409).

---

## 🛡️ Безопасность

### Защита эндпоинтов

Все эндпоинты админ-панели защищены двумя middleware:

1. **AuthMiddleware** - проверка JWT токена
2. **AdminMiddleware** - проверка роли администратора

### Логирование

Все административные действия логируются:
- Удаление пользователей
- Удаление объявлений
- Изменение ролей
- Одобрение / отклонение объявлений и профилей мастеров
- Управление чёрным списком

Пример лога:
```
INFO admin access granted user_id=1 email=admin@example.com
INFO user deleted by admin user_id=123
INFO user role updated by admin user_id=45 new_role=admin
INFO ad approved by admin ad_id=15
INFO ad rejected by admin ad_id=22
INFO worker profile approved by admin worker_id=12
INFO worker profile rejected by admin worker_id=7
```

---

## 📝 Примеры использования (curl)

### Вход как администратор
```bash
# 1. Логин
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'

# Сохраните токен из ответа
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Просмотр всех пользователей
```bash
curl -X GET "http://localhost:8080/admin/users?limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

### Изменение роли пользователя
```bash
curl -X PATCH http://localhost:8080/admin/users/5/role \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_name": "admin"
  }'
```

### Одобрение объявления
```bash
curl -X PATCH http://localhost:8080/admin/ads/15/approve \
  -H "Authorization: Bearer $TOKEN"
```

### Отклонение объявления
```bash
curl -X PATCH http://localhost:8080/admin/ads/22/reject \
  -H "Authorization: Bearer $TOKEN"
```

### Одобрение профиля мастера
```bash
curl -X PATCH http://localhost:8080/admin/workers/12/approve \
  -H "Authorization: Bearer $TOKEN"
```

### Удаление объявления
```bash
curl -X DELETE http://localhost:8080/admin/ads/15 \
  -H "Authorization: Bearer $TOKEN"
```

### Просмотр статистики
```bash
curl -X GET http://localhost:8080/admin/stats \
  -H "Authorization: Bearer $TOKEN"
```

### Добавление в черный список
```bash
curl -X POST http://localhost:8080/admin/blacklist \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "spam@example.com"
  }'
```

### Создание новой категории
```bash
curl -X POST http://localhost:8080/admin/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ландшафтный дизайн"
  }'
```

### Создание новой единицы цены
```bash
curl -X POST http://localhost:8080/admin/price-units \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "за м²"
  }'
```

---

## 🚀 Создание первого администратора

### Вариант 1: Через SQL
```sql
-- 1. Создать роль admin (если еще нет)
INSERT INTO roles (role_name, created_at, updated_at) 
VALUES ('admin', NOW(), NOW());

-- 2. Создать пользователя-администратора
INSERT INTO users (email, name, password_hash, role_id, created_at, updated_at)
VALUES (
  'admin@example.com',
  'Администратор',
  '$2a$10$...',  -- хеш пароля
  (SELECT id FROM roles WHERE role_name = 'admin'),
  NOW(),
  NOW()
);
```

### Вариант 2: Через API
```bash
# 1. Зарегистрироваться как обычный пользователь
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "name": "Администратор",
    "password": "secure_password",
    "role": 1
  }'

# 2. В БД вручную изменить role_id на роль admin
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE role_name = 'admin')
WHERE email = 'admin@example.com';
```

---

## ⚠️ Коды ошибок

| Код | Описание |
|-----|----------|
| 401 | Не авторизован (нет токена или токен невалиден) |
| 403 | Доступ запрещен (недостаточно прав, не администратор) |
| 404 | Ресурс не найден |
| 400 | Неверный формат запроса |
| 500 | Внутренняя ошибка сервера |

---

## 📊 Архитектура

```
/admin
├── /users          - Управление пользователями
│   ├── GET /       - Список пользователей
│   ├── GET /{id}   - Пользователь по ID
│   ├── DELETE /{id} - Удалить пользователя
│   └── PATCH /{id}/role - Изменить роль
│
├── /ads            - Модерация объявлений
│   ├── GET /           - Все объявления (?status=pending|approved|rejected)
│   ├── DELETE /{id}    - Удалить объявление
│   ├── PATCH /{id}/approve - Одобрить объявление
│   └── PATCH /{id}/reject  - Отклонить объявление
│
├── /workers        - Модерация профилей мастеров
│   ├── GET /           - Профили на модерации (?status=pending|approved|rejected)
│   ├── PATCH /{id}/approve - Одобрить профиль
│   └── PATCH /{id}/reject  - Отклонить профиль
│
├── /responses      - Модерация откликов
│   ├── GET /       - Все отклики
│   └── DELETE /{id} - Удалить отклик
│
├── /stats          - Статистика
│   └── GET /       - Общая статистика
│
├── /blacklist      - Черный список
│   ├── GET /       - Список
│   ├── POST /      - Добавить email
│   └── DELETE /{email} - Удалить email
│
├── /categories     - Управление категориями
│   ├── POST /      - Создать категорию
│   ├── PATCH /{id} - Обновить категорию
│   └── DELETE /{id} - Удалить категорию
│
└── /price-units    - Управление единицами цены
    ├── POST /      - Создать единицу цены
    ├── PATCH /{id} - Обновить единицу цены
    └── DELETE /{id} - Удалить единицу цены
```

---

## 🔧 Расширение функционала

Для добавления новых функций в админ-панель:

1. Добавьте хендлер в соответствующий файл (`users.go`, `moderation.go`)
2. Зарегистрируйте роут в `routes.go`
3. Все эндпоинты автоматически защищены `AdminMiddleware`

Пример:
```go
// В routes.go
admin.Get("/reports", GetReportsHandler(db, logger))

// В новом файле reports.go
func GetReportsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Ваша логика
    }
}
```
