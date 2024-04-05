include .env

docker-up:
	docker compose --profile $(ENV) up -d
docker-down:
	docker compose --profile $(ENV) down
docker-reset-down:
	docker compose --profile $(ENV) down -v
docker-bash:
	docker exec -it dcard_mysql bash
docker-mysql:
	docker exec -it dcard_mysql mysql -D dcard -u root -p
docker-redis:
	docker exec -it dcard_redis redis-cli
migrate-up:
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/dcard" -verbose up
migrate-down:
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/dcard" -verbose down