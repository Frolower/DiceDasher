<div align="center">
<h1>DiceDasher</h1>
  

  <a href="https://github.com/Frolower/DiceDasher/issues/new?assignees=&labels=bug&template=01_BUG_REPORT.md&title=bug%3A+">Report a Bug</a>
  ·
  <a href="https://github.com/Frolower/DiceDasher/issues/new?assignees=&labels=enhancement&template=02_FEATURE_REQUEST.md&title=feature%3A+">Request a Feature</a>
  ·
  <a href="https://github.com/Frolower/DiceDasher/issues/new?assignees=&labels=question&template=04_SUPPORT_QUESTION.md&title=support%3A+">Ask a Question</a>
</div>

<div align="center">

[![License](https://img.shields.io/github/license/Frolower/DiceDasher.svg?style=flat-square)](./LICENSE)

</div>

<details open="open">
<summary>Table of Contents</summary>

- [Status](#status)
- [About](#about)
  - [Built With](#built-with)
- [Architecture](#architecture)
  - [Repo Map](#repo-map)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Support](#support)
- [Project assistance](#project-assistance)
- [Contributing](#contributing)
- [Authors & contributors](#authors--contributors)
- [Security](#security)
- [License](#license)
- [Acknowledgements](#acknowledgements)

</details>

---
## Status
Early development. Backend APIs and a minimal React landing page with registration
and login are available. Open the frontend at http://localhost:5174. Dice rolls,
character creation, and profile management currently use the API or Swagger UI.
<br>
Currently supported systems:
- General
- VtM V5
- The Electric State

Currently existing services:
- backend service: handles registration, JWT authentication, session refresh, profile lookup/updates, password changes, and account deletion ([documentation](docs/services/backend/README.md))
- resolve service: resolves dice rolls according to the system rules
- character service: creates and stores system-specific characters

## About

DiceDasher is a website created to have a dedicated space to run online TTRPG games, with tools for specific systems.
It is designed to have both system agnostic and system specific tools for various TTRPG systems.

### Built With

- Backend: Go 1.26, net/http
- Frontend: React, TypeScript, Vite, Nginx
- DB: PostgreSQL 18 (UUID v7 defaults)
- Infra: Docker Compose, Makefile

## Architecture
This project is designed to use microservice architecture.

### Repo Map
```
.
├── database/                     
│   └── init/                     Database initialization scripts
│
├── docs/
│   ├── api/                      Existing API endpoints overview
│   ├── services/
│   │   ├── backend/             Authentication and user management documentation
│   │   ├── character/           Character service documentation
│   │   └── resolve/             Resolve service documentation
│   ├── CODE_OF_CONDUCT.md
│   ├── CONTRIBUTING.md
│   └── SECURITY.md
│
├── frontend/
│
├── pkg/                           Project-wide helper functions
│   ├── dbutil/                    DB related helpers
│   ├── dice/                      Dice rolling engine
│   ├── httputil/                  Web related helpers
│   ├── logger/                    Custom logging functions
│   ├── util/                      Various utility functions
│   └── go.mod
│
├── services/                      Backend services
│   ├── backend/                   Authentication, sessions, and user management
│   ├── character/                 Character service
│   │   ├── cmd/character-service  Entrypoint
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── handler/
│   │   │   ├── repository/
│   │   ├── .env.example
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   └── resolve/                   Resolve service
│       ├── cmd/resolve-service    Entrypoint
│       ├── internal/
│       │   ├── config/
│       │   ├── handler/
│       │   ├── repository/
│       │   └── system/            Per system resolvers (eg vtmv5/)
│       ├── .env.example
│       ├── Dockerfile
│       └── go.mod
│
├── .dockerignore
├── .gitignore
├── docker-compose.yml
├── go.work
├── LICENSE
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

The Docker setup requires Docker with Compose. The commands below also use Make;
without Make, run the corresponding `docker compose` commands directly. Services
run in containers, so local Go and Node.js installations are optional.

### Installation

1. Clone the repo
  ```sh
  git clone https://github.com/Frolower/DiceDasher
  ```
2. Enter the repository and create service environment files
  ```sh
  cd DiceDasher
  cp services/resolve/.env.example services/resolve/.env
  cp services/character/.env.example services/character/.env
  cp services/backend/.env.example services/backend/.env
  ```
  Set `JWT_SECRET` in `services/backend/.env` to the output of
  `openssl rand -base64 32`. Keep `COOKIE_SECURE=false` for local HTTP.
  Compose sets the frontend Origin to `http://localhost:5174`.
3. Build and run Docker containers
  ```sh
  make build
  make up
  ```
4. Display logs
  ```sh
  make logs
  ```
5. Check if services are healthy
  ```sh
  curl http://localhost:8080/health
  curl http://localhost:8081/health
  curl http://localhost:8083/health
  curl http://localhost:5174/health
  ```
6. Open the frontend at [http://localhost:5174](http://localhost:5174) and Swagger UI at [http://localhost:8082](http://localhost:8082).
7. Stop containers
  ```sh
  make down
  ```

### Make commands

Run `make` or `make help` to list all available commands.

| Command | Description |
| --- | --- |
| `make up` / `make down` | Start / stop the Docker stack |
| `make dev` | Build and run the Docker stack with attached logs |
| `make build` / `make rebuild` | Build images / rebuild without cache |
| `make logs SERVICE=backend` | Follow logs for a specific service |
| `make ps` / `make restart` | Check status / restart containers |
| `make install` | Install frontend dependencies from the lockfile |
| `make frontend-dev` | Start Vite with live reload |
| `make frontend-build` | Type-check and build the frontend |
| `make test` | Run Go and frontend tests |
| `make lint` | Run Go vet and frontend lint |
| `make fmt` / `make fmt-check` | Format files / check formatting |

Local Go commands require the Go version specified in `go.work`; frontend commands
require Node.js and npm. Run `make install` before local frontend commands.
Before `make frontend-dev`, stop the frontend container with
`docker compose stop frontend` and start the API with `make up SERVICE=backend`.
Both Vite and the frontend container use port 5174.

**Destructive commands:** `make clean` and `make nuke` delete database volumes and
all stored data. `make down` keeps volumes.

## Usage

The [frontend](frontend/README.md) provides registration and login. Its access token
is held in memory; reloading clears the displayed login state. Use Swagger UI,
Postman, or cURL for the other endpoints listed in the
[API reference](docs/api/README.md).

## Roadmap

See the [open issues](https://github.com/Frolower/DiceDasher/issues) for a list of proposed features (and known issues).

- [Top Feature Requests](https://github.com/Frolower/DiceDasher/issues?q=label%3Aenhancement+is%3Aopen+sort%3Areactions-%2B1-desc) (Add your votes using the 👍 reaction)
- [Top Bugs](https://github.com/Frolower/DiceDasher/issues?q=is%3Aissue+is%3Aopen+label%3Abug+sort%3Areactions-%2B1-desc) (Add your votes using the 👍 reaction)
- [Newest Bugs](https://github.com/Frolower/DiceDasher/issues?q=is%3Aopen+is%3Aissue+label%3Abug)

## Support

Reach out to the maintainer at one of the following places:

- [GitHub issues](https://github.com/Frolower/DiceDasher/issues/new?assignees=&labels=question&template=04_SUPPORT_QUESTION.md&title=support%3A+)
- Contact options listed on [this GitHub profile](https://github.com/Frolower)

## Project assistance

If you want to say **thank you** or/and support active development of DiceDasher:

- Add a [GitHub Star](https://github.com/Frolower/DiceDasher) to the project.
- Spread the word about this project among your friends and use it.

Together, we can make DiceDasher **better**!

## Contributing

First off, thanks for taking the time to contribute! Contributions are what make the open-source community such an amazing place to learn, inspire, and create.
Any contributions you make will benefit everybody else and are **greatly appreciated**.


Please read [our contribution guidelines](/docs/CONTRIBUTING.md), and thank you for being involved!

By participating, you agree to abide by the Code of Conduct.

## Authors & contributors

The original setup of this repository is by [Nikita Dubinin](https://github.com/Frolower).

For a full list of all authors and contributors, see [the contributors page](https://github.com/Frolower/DiceDasher/contributors).

## Security

DiceDasher follows good practices of security, but 100% security cannot be assured.
DiceDasher is provided **"as is"** without any **warranty**. Use at your own risk.

_For more information and to report security issues, please refer to our [security documentation](/docs/SECURITY.md)._

## License

This project is licensed under the **MIT license**.

See [LICENSE](/LICENSE) for more information.

## Acknowledgements

> This project is inspired by the TTRPGs, platforms like owlbear.rodeo and Roll20
