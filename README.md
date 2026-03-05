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
Early development. Backend API is available; UI is not implemented yet. 
Use Postman/cURL for now.
<br>
Currently supported systems:
- General
- VtM V5
- The Electric State

Currently existing services:
- resolve service: resolves dice rolls according to the system rules

## About

DiceDasher is a website created to have a dedicated space to run online TTRPG games, with tools for specific systems.
It is designed to have both system agnostic and system specific tools for various TTRPG systems.

### Built With

- Backend: Go 1.25, net/http
- Frontend: React
- DB: PostgreSQL
- Infra: Docker Compose, Makefile

## Architecture
This project is designed to use microservice architecture.

### Repo Map
```
.
├── database/                     Database initialization scripts
│
├── docs/
│   ├── service/
│   │   └──resolve/               Documentation about resolve service
│   ├── CODE_OF_CONDUCT.md
│   ├── CONTRIBUTING.md
│   ├── ENDPOINTS.md              Implemented API endpoints
│   └── SECURITY.md
│
├── frontend/
│
├── pkg/                           Project-wide helper functions
│   ├── dice/                      Dice rolling engine
│   ├── httputil/                  Web related helpers
│   ├── logger/                    Custom logging functions
│   ├── util/                      Various utility functions
│   └── go.mod
│
├── services/                       Backend services
│   ├── api/                        Existing API endpoints overview
│   └── resolve-service/            Resolve service
│       ├── cmd/resolve-service     Entrypoint
│       ├── internal/
│       │   ├── config/
│       │   ├── handler/
│       │   ├── repository/
│       │   ├── system/             Per system resolvers (eg vtmv5/)
│       ├── Dockerfile
│       └── go.mod
│
├── .gitignore
├── docker-compose.yml
├── go.work
├── LICENSE
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

The only current requirement is the latest release of Docker, all services are packed
into docker containers and won't require any additional downloads

### Installation

1. Clone the repo
  ```sh
  git clone https://github.com/Frolower/DiceDasher
  ```
2. Run docker containers
  ```sh
  make up
  ```
3. Display logs
  ```sh
  make logs
  ```
4. Check if services are healthy
  ```sh
  curl http://localhost:8080/health
  ```
5. Stop containers
  ```sh
  make down
  ```

## Usage

So far there is no user interface to access project functionality. The recommended
way is to use Postman to access API endpoints. All currently supported requests are
listed in the [documentation](/docs/api/README.md).

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
