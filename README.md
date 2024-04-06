# 2024 Backend Intern Assignment (Dcard)

## About The Project

### Built With
- Golang 1.22.1 (gin, go-redis, sqlc, go-swagger)
- MySQL 8.0.36
- Redis 7.2.4
- k6 v0.50.0

## Getting Started

To run the app, follow these steps:

1. Add your application configuration to `.env` file in the root of the project:

```sh
# at .env
ENV=prod

MYSQL_ROOT_PASSWORD=
MYSQL_DATABASE=
MYSQL_USER=
MYSQL_PASSWORD=
```

2. Build and run the app using Docker Compose:

```
docker compose up --build
```

3. Access the API endpoint at `localhost:8080/api/v1/ad`

4. For Swagger API documentation, visit `localhost:8080/api/v1/swagger/index.html`

## Database Design

![database design](docs/database_design.png)

## Load Testing with Grafana k6

```
cd k6
python pre.py # create fake advertisements
k6 run script.js
```

![k6 result](docs/k6_result.png)
