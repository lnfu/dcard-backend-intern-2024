# 2024 Backend Intern Assignment (Dcard)

## About The Project

### Built With
- Golang 1.22.1
- MySQL 8.0.36
- Redis 7.2.4
- k6 v0.50.0

## Getting Started

Add your application configuration to `.env` file in the root of the project:

```sh
# at .env
ENV=dev

MYSQL_ROOT_PASSWORD=
MYSQL_DATABASE=
MYSQL_USER=
MYSQL_PASSWORD=
```

```
docker compose up --build
```

![database design](docs/database_design.png)

![k6 result](docs/k6_result.png)
