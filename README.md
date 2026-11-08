# DiceDasher

DiceDasher is a platform in development for online tabletop role-playing games,
with system-independent tools and rules for individual TTRPGs.

[Report a bug](https://github.com/Frolower/DiceDasher/issues) ·
[Request a feature](https://github.com/Frolower/DiceDasher/issues) ·
[License](LICENSE)

## Current state

The web app includes registration, login, session restoration after a page reload,
and account settings. Signed-in users can edit their username and email, change
their password and delete their account. Access tokens stay in memory;
sessions are restored and renewed using an HttpOnly refresh cookie.

Rooms, Characters, and Friends screens currently show example content. Room
creation, character creation in the UI, and friend requests are placeholders.
Systems mentioned on example cards do not indicate implemented rules support.
Dice resolution and character creation are available through separate APIs and
Swagger UI, rather than the web interface.

| Service | Implemented functionality |
| --- | --- |
| `backend` | Registration, JWT authentication, rotating refresh tokens, logout, profile lookup/updates, password changes, account deletion |
| `resolve` | Generic dice rolls; Vampire: The Masquerade V5 rolls, rerolls, and checks; The Electric State rolls and pushes; persisted roll history |
| `character` | Validation, creation, and storage of The Electric State characters |

## Stack and architecture

- **Backend:** Go, `net/http`, three services in a Go workspace. The workspace currently specifies Go **1.27.2**.
- **Frontend:** React 19, TypeScript, React Router, Vite; static builds served by Nginx.
- **Database:** PostgreSQL 18 with pg_cron, separate databases and roles for each service, UUID v7 identifiers.
- **Development:** Docker Compose, Make, Node.js 24, npm, Oxlint, Prettier, GitHub Actions.

The browser sends account requests to `/api/*`. Nginx (or Vite in local
development) forwards them to `backend`, removing the `/api` prefix. The resolve
and character APIs are accessed separately; the frontend proxy does not route to
them. Each Go service connects to its own database in the shared PostgreSQL
instance.

```text
.
├── database/                 PostgreSQL image, initialization SQL, database tests
├── docs/
│   ├── api/                  API overview and OpenAPI specification
│   └── services/             Backend, resolve, and character documentation
├── frontend/                 React app, Nginx configuration, frontend checks
├── pkg/                      Shared dice, database, HTTP, logging, and utility code
├── services/
│   ├── backend/              Authentication, sessions, and account management
│   ├── character/            Character validation and persistence
│   └── resolve/              System-specific dice resolution and roll history
├── docker-compose.yml        Shared service configuration
├── docker-compose.dev.yml    Local ports, HTTP settings, and Swagger UI
├── go.work                   Go workspace
└── Makefile                  Development commands
```

## Quick start

Install Docker with Docker Compose. The commands below use Make; local Go and
Node.js installations are not needed for the container setup.

```sh
git clone https://github.com/Frolower/DiceDasher.git
cd DiceDasher
cp services/backend/.env.example services/backend/.env
cp services/resolve/.env.example services/resolve/.env
cp services/character/.env.example services/character/.env
```

Configure the new `.env` files before starting:

| File | Required setting for Docker Compose |
| --- | --- |
| `services/backend/.env` | Set `JWT_SECRET` to the output of `openssl rand -base64 32` |
| `services/resolve/.env` | Set `DATABASE_URL=postgres://resolve_service:resolve_password@postgres:5432/resolve_db?sslmode=disable` |
| `services/character/.env` | Set `DATABASE_URL=postgres://character_service:character_password@postgres:5432/character_db?sslmode=disable` |

Compose supplies the backend database URL. Leave the example HTTP ports unchanged;
container health checks use them. The development override sets
`COOKIE_SECURE=false` and allows browser requests to the backend from
`http://localhost:5174` and `http://localhost:8888`.

```sh
make up
make ps
```

`make up` builds images and starts the stack in the background. On the first
startup, PostgreSQL creates the databases and loads the initialization and seed
scripts.

| Component | Local address |
| --- | --- |
| Web app | http://localhost:5174 |
| Swagger UI | http://localhost:8888 |
| Account API | http://localhost:8000 |
| Resolve API | http://localhost:8081 |
| Character API | http://localhost:8082 |
| PostgreSQL | `localhost:5432` |

All published development ports bind to `127.0.0.1`.

```sh
# Check HTTP health endpoints
curl http://localhost:5174/health
curl http://localhost:8000/health
curl http://localhost:8081/health
curl http://localhost:8082/health

# Follow logs, or stop the stack while keeping database data
make logs
make down
```

### Compose configuration and data

Docker Make targets explicitly combine `docker-compose.yml` and
`docker-compose.dev.yml`. The equivalent startup command is:

```sh
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
```

Use both files for other Compose operations too. A plain `docker compose up`
loads only the base configuration and publishes no host ports. To inspect the
merged configuration:

```sh
docker compose -f docker-compose.yml -f docker-compose.dev.yml config
```

The stack is configured for local development, with development database
credentials and HTTP cookies. No production Compose or TLS configuration is
provided. `COMPOSE_FILES` can override the file selection used by Make.

PostgreSQL data lives in the `postgres_data` volume. Initialization SQL runs only
on an empty volume; changing a script does not migrate an existing database.
**`make clean` and `make nuke` delete database volumes and stored data.**
See the [database documentation](database/README.md) for schema updates,
timestamp triggers, session retention, and pg_cron checks.

### Local frontend development

Use Node.js 24 and npm, matching the frontend image and CI. From the repository
root, after configuring the service environment files:

```sh
make install
docker compose -f docker-compose.yml -f docker-compose.dev.yml stop frontend
make up SERVICE=backend
make frontend-dev
```

Open http://localhost:5174. Vite provides live reload and proxies `/api/*` to
`http://localhost:8000`. Vite and the frontend container share port 5174, so stop
the container before starting Vite. Rebuild the frontend container with
`make up SERVICE=frontend` when returning to the Docker-served app.

Local Go commands require the version specified in [go.work](go.work). When
running Go services outside Docker, supply their environment variables explicitly
and use `localhost` instead of `postgres` in database URLs.

## Web routes

| Route | Behavior |
| --- | --- |
| `/` | Public landing page; signed-in users redirect to `/profile` |
| `/login` | Login dialog; returns to the requested private page after login |
| `/profile` | Signed-in profile with example rooms |
| `/rooms` | Redirects to `/profile` |
| `/settings` | Account details, password, and account deletion |
| `/characters` | Signed-in character preview with example data |
| `/friends` | Signed-in friends preview with example data |

Private routes wait for session restoration before redirecting to login.
Nginx serves the app for direct links to client-side routes.

## API documentation

Open [Swagger UI](http://localhost:8888) with the development stack running, or
read the [API overview](docs/api/README.md) and
[OpenAPI specification](docs/api/openapi.yaml).

- [Account API](docs/services/backend/README.md)
- [Resolve endpoints and supported actions](docs/services/resolve/endpoints.md)
- [Generic dice](docs/services/resolve/systems/generic.md), [VtM V5](docs/services/resolve/systems/vtmv5.md), [The Electric State](docs/services/resolve/systems/tes.md)
- [Character endpoints](docs/services/character/endpoints.md) and [TES character schema](docs/services/character/systems/tes.md)

Account endpoints use paths such as `/login` and `/me` on port 8000, or
`/api/login` and `/api/me` through the frontend. Resolve requests use
`POST /resolve?system=generic|vtmv5|tes`; character creation uses
`POST /character?system=tes` on its own service.

The Docker frontend applies Nginx request limits to `/api/`, with additional
limits for login and registration. Direct requests to port 8000 and requests
through Vite bypass these Nginx limits. Shared Nginx routes and limits live in
`frontend/nginx-common/`; the dev Compose file mounts this directory and
`frontend/nginx.dev.conf`. After editing mounted configuration, validate and reload:

```sh
docker compose -f docker-compose.yml -f docker-compose.dev.yml exec frontend nginx -t
docker compose -f docker-compose.yml -f docker-compose.dev.yml exec frontend nginx -s reload
```

## Development commands and checks

Run `make help` for the full list. Frontend commands require `make install` first;
Go commands cover all four workspace modules.

| Command | Purpose |
| --- | --- |
| `make up` / `make down` | Build and start the dev stack / stop it and keep data |
| `make dev` | Start the stack with attached logs |
| `make build` / `make rebuild` | Build images / build without cache |
| `make logs SERVICE=backend` | Follow one service's logs |
| `make ps` / `make restart` | Show status / restart containers |
| `make frontend-dev` | Start Vite with live reload |
| `make frontend-build` | Type-check and build the frontend |
| `make frontend-preview` | Preview the built frontend |
| `make test` | Run Go tests and frontend authentication/settings checks |
| `make lint` | Run Go vet and Oxlint |
| `make fmt` / `make fmt-check` | Format code / check formatting |

`SERVICE` selects a service for build, up, dev, logs, restart, and rebuild.
For separate checks, use `make test-go`, `make test-frontend`, `make lint-go`, or
`make lint-frontend`.

Pull requests to `master` run frontend lint/tests and image build, Go vet and tests
with the race detector, database SQL tests, backend PostgreSQL integration tests,
and Go service image builds. `make test` does not provision the integration-test
database; see [backend testing instructions](docs/services/backend/README.md) and
[database checks](database/README.md) for database-backed verification.

## Contributing and support

See [open issues](https://github.com/Frolower/DiceDasher/issues) for bugs, feature
requests, and development discussion. Read the
[contribution guidelines](docs/CONTRIBUTING.md) and
[Code of Conduct](docs/CODE_OF_CONDUCT.md) before contributing.
Report security issues according to the [security policy](docs/SECURITY.md).

Created by [Nikita Dubinin](https://github.com/Frolower).
See all [contributors](https://github.com/Frolower/DiceDasher/contributors).
Inspired by tabletop RPGs and platforms such as Owlbear Rodeo and Roll20.

## License

Licensed under the [Apache License 2.0](LICENSE).
