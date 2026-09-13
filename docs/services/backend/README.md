# Backend: регистрация и вход

Backend: `http://localhost:8083`. Регистрация не создаёт сессию. Вход по username
и паролю создаёт отдельную сессию и выдаёт access JWT + refresh-cookie.

## Где находится логика

| Файл / пакет | Ответственность |
|---|---|
| `cmd/backend-service/main.go` | Сборка repository, bcrypt, JWTManager, auth.Service и handler |
| `internal/handler/handler.go` | HTTP-интерфейс Authenticator, зависимости и настройки cookie |
| `internal/handler/register.go` | HTTP-контракт регистрации |
| `internal/handler/login.go` | JSON входа, проверка Origin, ошибки и ответ с токеном |
| `internal/handler/cookie.go` | Refresh-cookie |
| `internal/handler/cors.go` | CORS для конкретного Origin с credentials |
| `internal/auth/service.go` | Интерфейсы хранилищ и конструкторы |
| `internal/auth/register.go` | Регистрация и сохранение bcrypt-хеша |
| `internal/auth/login.go` | Поиск пользователя и проверка пароля |
| `internal/auth/session.go` | Создание сессии после проверки пароля |
| `internal/auth/jwt.go` | Выпуск и проверка access JWT |
| `internal/auth/refresh_token.go` | Криптографически случайный refresh и SHA-256 |
| `internal/auth/errors.go`, `model.go`, `validate.go` | Ошибки, модели и валидация |
| `internal/repository/user.go`, `session.go` | SQL пользователей и сессий |

HTTP-слой не выполняет SQL и bcrypt. Сервис не устанавливает cookie. Репозиторий
не выпускает токены. Context переносит отмену запроса и данные логирования.

## POST /register

Вход: `username`, `email`, `password`. Успех: `201 {"status":"created"}` без cookie.

- Username: 5–64 Unicode-символа, без пробельных и управляющих символов; не обрезается.
- Email: 3–254 байта, проверка формата, trim и нижний регистр.
- Пароль: минимум 8 символов, максимум 72 байта UTF-8; строчные, заглавные буквы и цифра; без пробельных символов. Не изменяется.
- Username и email уникальны без учёта регистра через PostgreSQL `lower(...)`.
- Ошибки: 400 — ввод, 409 — дубликат, 413 — тело больше 16 КиБ, 500 — внутренний сбой.

## POST /login

```sh
curl -i -c /tmp/dicedasher-cookies.txt http://localhost:8083/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"Alice","password":"StrongPass1"}'
```

Успех: HTTP `200`, `Cache-Control: no-store`:

```json
{
  "access_token": "<signed JWT>",
  "token_type": "Bearer",
  "expires_at": "<UTC timestamp>"
}
```

Refresh-токен передаётся только через `Set-Cookie`, а не JSON:

- HTTPS: `__Host-refresh_token`, Secure, HttpOnly, SameSite=Strict, Path=/, без Domain.
- Локальный HTTP с `COOKIE_SECURE=false`: `refresh_token`, HttpOnly, SameSite=Strict, Path=/.
- Expires и Max-Age ограничены сроком сессии.

В браузере запрос должен иметь `credentials: "include"`, в том числе при первом
входе. CORS разрешает только `ALLOWED_ORIGIN`. Login принимает только
`application/json`, отвергает чужой Origin и межсайтовый Fetch Metadata без Origin.
Запросы curl/Postman без Origin допустимы. User-Agent берётся из HTTP-заголовка;
он используется для отображения, а не доказательства личности устройства.

| Статус | Причина |
|---|---|
| 200 | Пароль проверен, сессия сохранена, токены выданы |
| 400 | Некорректный JSON, отсутствующие/слишком длинные поля |
| 401 | Неверный username или пароль (одинаковый ответ) |
| 403 | Неразрешённый Origin |
| 413 | Тело больше 16 КиБ |
| 415 | Content-Type не application/json |
| 500 | Внутренняя ошибка; токены и cookie не выдаются |

При входе не проверяется новая сложность пароля регистрации: старые пароли
должны продолжать работать. Пароль проверяется через bcrypt.CompareHashAndPassword,
а не повторным хешированием. Для отсутствующего пользователя также выполняется
bcrypt-проверка dummy-хеша, чтобы уменьшить различие времени ответа.

## Токены и сессии

Access JWT (HS256) действует 10 минут. Содержит sub, sid, iss, aud, iat, exp и jti.
Проверка ограничена HS256, ожидаемыми issuer/audience и обязательными exp/iat;
валидируются UUID пользователя и сессии. JWT подписан, но его содержимое не зашифровано.

Refresh — 32 случайных байта, закодированных base64url. В БД сохраняется только
SHA-256-хеш, не исходный токен. Сессия действует 30 дней. Каждый вход создаёт новую
сессию. Идентификатор сессии генерируется сервисом до подписи JWT и INSERT.

Сначала подготавливаются токены, затем записывается сессия, затем отправляется ответ.
Ошибка подписи не создаёт строку БД; ошибка записи не выдаёт токены. Обрыв сети после
успешного INSERT может оставить сессию, которую клиент не получил; login не идемпотентен.

Refresh и logout пока возвращают ErrNotImplemented на уровне сервиса; публичный
logout остаётся 501, маршрут refresh ещё не подключён. Методы SQL для ротации/отзыва
подготовлены, но история использованных refresh-токенов и защита от их повторного
использования ещё не реализованы. После истечения access JWT пока нужен новый login.
Проверка JWT реализована как компонент, но middleware защиты `/me` и других
маршрутов — следующий этап. Подписанный JWT сам не проверяет revoked_at в БД.

## Настройка и запуск

Нужны значения из `services/backend/.env.example`. Создайте локальный `.env` и один
раз запишите в JWT_SECRET результат `openssl rand -base64 32`. Не коммитьте ключ.
Если локальный ключ уже создан, сохраняйте его между перезапусками; смена ключа
делает старые access JWT недействительными. Пустой/короткий ключ останавливает старт.

| Переменная | Значение |
|---|---|
| DATABASE_URL | Обязательный URL user_db |
| HTTP_ADDR | По умолчанию :8083 |
| JWT_SECRET | Обязательный base64-секрет минимум из 32 случайных байт |
| JWT_ISSUER | dicedasher-backend |
| JWT_AUDIENCE | dicedasher-api |
| COOKIE_SECURE | По умолчанию true; false только для локального HTTP |
| ALLOWED_ORIGIN | http://localhost:8082, без завершающего слеша |
| LOG_MODE | default |

Compose читает `services/backend/.env`, но переопределяет адрес БД на внутренний
`postgres:5432/user_db` с ролью user_service. Прямой запуск Go сам .env не читает.

Для существующего PostgreSQL volume:

```sh
docker compose up -d postgres
docker compose exec -T postgres psql -U postgres -d postgres < database/init/004_backend.sql
docker compose up -d --build backend
```

004 можно повторно применить: он создаёт отсутствующие user_db/роль, таблицы,
индексы и права. Не меняет пароль существующей роли и не мигрирует произвольную
старую структуру таблиц. Данные из прежней backend_db автоматически не переносятся.
Новый volume инициализируется через 000 и 004 автоматически. Не удаляйте volume
ради применения новых GRANT/индексов.

В рабочем окружении используйте HTTPS, COOKIE_SECURE=true и свой Origin.
HS256-ключ даёт право и проверять, и выпускать JWT; не передавайте его клиенту.
Ограничение частоты попыток входа пока не добавлено; перед публичным запуском
нужен rate limit на ingress или в приложении.

## Проверки

```sh
go test ./services/backend/...
go vet ./services/backend/...
```

Для интеграционных тестов примените 004 к отдельной тестовой БД и укажите:

```sh
BACKEND_TEST_DATABASE_URL='postgres://user_service:user_password@localhost:55439/user_db?sslmode=disable' \
BACKEND_TEST_ADMIN_URL='postgres://YOUR_ADMIN@localhost:55439/user_db?sslmode=disable' \
go test -race ./services/backend/...
```

Тесты проверяют регистрацию, bcrypt, JWT, cookie, повторный вход, ошибки пароля,
конкурирующую регистрацию и реальные SQL-права роли. Тестовые пользователи имеют
случайные имена и удаляются вместе со своими сессиями.
