# Backend: регистрация пользователей

Backend слушает `http://localhost:8083`. Регистрация создаёт пользователя, но не
создаёт сессию, cookie или токен. Вход и остальные пользовательские маршруты
пока возвращают `501 Not Implemented`.

## Архитектура

- `cmd/backend-service/main.go` создаёт пул БД, repository, auth.Service и handler.
- `internal/handler` отвечает за JSON, ограничения тела и HTTP-статусы. Зависит от интерфейса `Registrar`.
- `internal/auth` выполняет регистрацию: нормализация → валидация → bcrypt → сохранение. Зависит от `Store` и `PasswordHasher`, не от HTTP или PostgreSQL.
- `internal/repository` выполняет параметризованный INSERT и переводит ошибки уникальности в `auth.ErrAlreadyExists`.

Зависимости передаются конструкторам, а context используется для отмены запроса
и логирования. Успех возвращается только после записи в БД.

## POST /register

```sh
curl -i http://localhost:8083/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"Alice","email":"alice@example.com","password":"StrongPass1"}'
```

Ответ: HTTP `201 Created`, `Content-Type: application/json`:

```json
{"status":"created"}
```

Правила:

- Все три поля обязательны; неизвестные JSON-поля отклоняются.
- Username: 5–64 Unicode-символа. Пробельные символы запрещены, в том числе по краям. Регистр сохраняется для отображения; управляющие символы запрещены.
- Email: 3–254 байта, проверка формата; пробелы по краям удаляются, регистр приводится к нижнему.
- Пароль: минимум 8 Unicode-символов, максимум 72 байта UTF-8; нужны строчная буква, заглавная буква и цифра. Пробельные символы запрещены, в том числе по краям; регистр сохраняется.
- Username и email уникальны без учёта регистра через индексы PostgreSQL `lower(...)` (правила регистра зависят от локали БД).
- Максимальный размер тела — 16 КиБ. Пароль хранится только как bcrypt-хеш.

| Статус | Значение |
|---|---|
| 201 | Пользователь сохранён |
| 400 | Ошибка JSON или валидации |
| 409 | Username или email уже занят |
| 413 | Слишком большое тело запроса |
| 500 | Ошибка хеширования или сохранения |

Ошибки возвращаются как `text/plain`; внутренние детали БД клиенту не передаются.
`GET /health` возвращает `200 OK` и текст `OK`; это проверка HTTP-сервиса, не постоянная проверка БД.

## Запуск

Новая локальная БД:

```sh
docker compose up -d --build backend
```

Если PostgreSQL volume уже существовал, init-скрипты автоматически повторно не
выполняются. Примените добавленный скрипт без удаления данных:

```sh
docker compose up -d postgres
docker compose exec -T postgres psql -U postgres -d postgres < database/init/004_backend.sql
docker compose up -d --build backend
```

Скрипт можно повторно запускать. Он создаёт `backend_db`, локальную роль
`backend_service` и таблицу `public.users`. Он не изменяет структуру ранее
созданной вручную таблицы users и не меняет пароль существующей роли.
Пароль в Compose/SQL — настройка локальной разработки.

Для запуска Go на хосте используйте значения из `services/backend/.env.example`
как переменные окружения (приложение само `.env` не читает).
`DATABASE_URL` обязателен, `HTTP_ADDR` по умолчанию `:8083`, `LOG_MODE` — `default`.
CORS разрешает Swagger UI с `http://localhost:8082`.

## Проверки

Из корня проекта с настроенным Go workspace:

```sh
go test ./services/backend/...
```

HTTP и сервисные тесты работают без БД. Для интеграционного теста примените
`004_backend.sql` в отдельном тестовом PostgreSQL и укажите два URL к backend_db:

```sh
BACKEND_TEST_DATABASE_URL='postgres://backend_service:backend_password@localhost:55439/backend_db?sslmode=disable' \
BACKEND_TEST_ADMIN_URL='postgres://YOUR_ADMIN@localhost:55439/backend_db?sslmode=disable' \
go test ./services/backend/internal/handler -run TestRegistrationPostgres -v
```

Тест проверяет реальную запись, bcrypt, конфликты email/username и конкурирующие
регистрации. Удаляет только собственные записи со случайными именами.
