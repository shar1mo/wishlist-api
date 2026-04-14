# Wishlist API

REST API сервис для создания вишлистов к праздникам и событиям.

Пользователь может:
- зарегистрироваться и войти по email и паролю
- создавать, просматривать, обновлять и удалять свои вишлисты
- добавлять, просматривать, обновлять и удалять подарки внутри вишлистов
- делиться вишлистом по публичной ссылке
- позволять другим пользователям бронировать подарки без авторизации по публичному токену

## Стек

- Go 1.24
- Chi
- PostgreSQL
- pgx
- Docker/Docker Compose
- golang-migrate
- JWT
- bcrypt

## Возможности

### Авторизация
- регистрация по email и паролю
- логин по email и паролю, получение jwt токена
- пароли хранятся в хэшированном виде
- закрытые эндпоинты доступны только с JWT токеном

### Вишлисты
- создание вишлиста
- получение списка своих вишлистов
- получение одного своего вишлиста
- обновление своего вишлиста
- удаление своего вишлиста

### Подарки
- создание подарка внутри вишлиста
- получение списка подарков вишлиста
- получение одного подарка
- обновление подарка
- удаление подарка

### Публичный доступ
- просмотр вишлиста по публичному токену
- бронирование подарка по публичному токену без авторизации
- защита от двойного бронирования через атомарный SQL-запрос

## Запуск

Проект запускается одной командой:

```bash
docker-compose up --build
```
если используется новая версия docker CLI:
```bash
docker compose up --build
```
после запуска API будет доступен по адресу:
```bash
http://localhost:8080
```

## Конфиг
Все параметры вынесены в переменные окружения.
В репо есть `.env.example` со списком доступных переменных окружения.

По дефолту проект может запускаться и без `.env`, так как в `docker-compose.yml` заданы значения по умолчанию.

## Архитектура проекта

```bash
.
├── cmd/app                 # точка входа
├── internal/auth           # JWT и хэширование паролей
├── internal/config         # конфиг
├── internal/handler        # HTTP handlers
├── internal/middleware     # middleware
├── internal/model          # модели
├── internal/repository     # интерфейсы и postgres-реализации
├── internal/service        # бизнес-логика
├── migrations              # SQL миграции
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API

### HTTP коды
- 200 OK
- 201 Created
- 204 No Content
- 400 Bad Request
- 401 Unauthorized
- 404 Not Found
- 409 Conflict
- 422 Unprocessable Entity
- 500 Internal Server Error

### Base URL for API endpoints:
```
http://localhost:8080/api/v1
```

### Healthcheck

```
http://localhost:8080/health
```

response:
```json
{
  "status": "ok"
}
```

### Auth
- Регистрация пользователя
```
POST /auth/register
```

request body:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

response `201 Created`:
```json
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2026-04-14T15:00:00Z"
}
```
- Логин пользователя
```
POST /auth/login
```
request body:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

response `200 OK`:
```json
{
  "token": "jwt-token"
}
```

### Protected endpoints

- для всех закрытых эндпоинтов требуется header:
```
Authorization: Bearer <jwt-token>
```

### Wishlists
- Создать вишлист
```
POST /wishlists
```
request body:
```json
{
  "title": "Birthday 2026",
  "description": "My birthday wishlist",
  "event_date": "2026-07-10"
}
```

response `201 Created`:
```json
{
  "id": 1,
  "title": "Birthday 2026",
  "description": "My birthday wishlist",
  "event_date": "2026-07-10",
  "public_token": "generated-public-token",
  "created_at": "2026-04-14T15:00:00Z",
  "updated_at": "2026-04-14T15:00:00Z"
}
```

- Получить список своих вишлистов
```
GET /wishlists
```

response `200 OK`:
```json
{
  "wishlists": [
    {
      "id": 1,
      "title": "Birthday 2026",
      "description": "My birthday wishlist",
      "event_date": "2026-07-10",
      "public_token": "generated-public-token",
      "created_at": "2026-04-14T15:00:00Z",
      "updated_at": "2026-04-14T15:00:00Z"
    }
  ]
}
```

- Получить вишлист по id
```
GET /wishlists/{wishlistId}
```
response `200 OK`:
```json
{
  "id": 1,
  "title": "Birthday 2026",
  "description": "My birthday wishlist",
  "event_date": "2026-07-10",
  "public_token": "generated-public-token",
  "created_at": "2026-04-14T15:00:00Z",
  "updated_at": "2026-04-14T15:00:00Z"
}
```

- Обновить вишлист
```
PUT /wishlists/{wishlistId}
```

request body:
```json
{
  "title": "Updated birthday wishlist",
  "description": "Updated description",
  "event_date": "2026-07-15"
}
```

response `200 OK`:
```json
{
  "id": 1,
  "title": "Updated birthday wishlist",
  "description": "Updated description",
  "event_date": "2026-07-15",
  "public_token": "generated-public-token",
  "created_at": "2026-04-14T15:00:00Z",
  "updated_at": "2026-04-14T15:10:00Z"
}
```

- Удалить вишлист
```
DELETE /wishlists/{wishlistId}
```
response `204 No Content`

### Items
- Создать подарок внутри вишлиста
```
POST /wishlists/{wishlistId}/items
```

request body:
```json
{
  "title": "Keyboard",
  "description": "logitech k380",
  "product_url": "https://example.com/keyboard",
  "priority": 5
}
```

response `201 Created`:
```json
{
  "id": 1,
  "wishlist_id": 1,
  "title": "Keyboard",
  "description": "logitech k380",
  "product_url": "https://example.com/keyboard",
  "priority": 5,
  "is_reserved": false,
  "reserved_at": null,
  "created_at": "2026-04-14T15:20:00Z",
  "updated_at": "2026-04-14T15:20:00Z"
}
```

- Получить список подарков вишлиста
```
GET /wishlists/{wishlistId}/items
```

response `200 OK`:
```json
{
  "items": [
    {
      "id": 1,
      "wishlist_id": 1,
      "title": "Keyboard",
      "description": "logitech k380",
      "product_url": "https://example.com/keyboard",
      "priority": 5,
      "is_reserved": false,
      "reserved_at": null,
      "created_at": "2026-04-14T15:20:00Z",
      "updated_at": "2026-04-14T15:20:00Z"
    }
  ]
}
```

- Получить подарок по id
```
GET /wishlists/{wishlistId}/items/{itemId}
```

response `200 OK`:
```json
{
  "id": 1,
  "wishlist_id": 1,
  "title": "Keyboard",
  "description": "logitech k380",
  "product_url": "https://example.com/keyboard",
  "priority": 5,
  "is_reserved": false,
  "reserved_at": null,
  "created_at": "2026-04-14T15:20:00Z",
  "updated_at": "2026-04-14T15:20:00Z"
}
```

- Обновить подарок
```
PUT /wishlists/{wishlistId}/items/{itemId}
```

request body:
```json
{
  "title": "Keyboard V2",
  "description": "Updated description",
  "product_url": "https://example.com/keyboard-v2",
  "priority": 4
}
```

response `200 OK`:
```json
{
  "id": 1,
  "wishlist_id": 1,
  "title": "Keyboard V2",
  "description": "Updated description",
  "product_url": "https://example.com/keyboard-v2",
  "priority": 4,
  "is_reserved": false,
  "reserved_at": null,
  "created_at": "2026-04-14T15:20:00Z",
  "updated_at": "2026-04-14T15:30:00Z"
}
```

- Удалить подарок
```
DELETE /wishlists/{wishlistId}/items/{itemId}
```
response `204 No Content`

### Public endpoints
- Публичные эндпоинты не требуют авторизации
- Получить публичный вишлист по токену
```
GET /public/wishlists/{token}
```

response `200 OK`:
```json
{
  "id": 1,
  "title": "Birthday 2026",
  "description": "My birthday wishlist",
  "event_date": "2026-07-10",
  "items": [
    {
      "id": 1,
      "title": "Keyboard",
      "description": "logitech k380",
      "product_url": "https://example.com/keyboard",
      "priority": 5,
      "is_reserved": false
    }
  ]
}
```

- Забронировать подарок по публичному токену
```
POST /public/wishlists/{token}/items/{itemId}/reserve
```

response `200 OK`:
```json
{
  "id": 1,
  "is_reserved": true,
  "reserved_at": "2026-04-14T15:40:00Z"
}
```

Если подарок уже забронирован, сервис возвращает:
response `409 Conflict`:
```json
{
  "error": "item already reserved"
}
```

## Примеры запросов

### Регистрация
```
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"user@example.com",
    "password":"password123"
  }'
```

### Логин
```
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"user@example.com",
    "password":"password123"
  }'
```

### Создать вишлист
```
curl -X POST http://localhost:8080/api/v1/wishlists \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Birthday 2026",
    "description":"My birthday wishlist",
    "event_date":"2026-07-10"
  }'
```

### Создать подарок
```
curl -X POST http://localhost:8080/api/v1/wishlists/1/items \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Keyboard",
    "description":"logitech k380",
    "product_url":"https://example.com/keyboard",
    "priority":5
  }'
```

### Получить публичный вишлист
```
curl http://localhost:8080/api/v1/public/wishlists/<token>
```

### Забронировать подарок
```
curl -X POST http://localhost:8080/api/v1/public/wishlists/<token>/items/1/reserve
```