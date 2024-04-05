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
	docker exec -it dcard_mysql mysql -D $(MYSQL_DATABASE) -u root -p
docker-redis:
	docker exec -it dcard_redis redis-cli
migrate-up:
	migrate -path db/migrations -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:3306)/$(MYSQL_DATABASE)" -verbose up
migrate-down:
	migrate -path db/migrations -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(localhost:3306)/$(MYSQL_DATABASE)" -verbose down
