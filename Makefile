include .env

docker-up:
	sudo docker compose --profile $(ENV) up -d
docker-down:
	sudo docker compose --profile $(ENV) down
docker-reset-down:
	sudo docker compose --profile $(ENV) down -v
docker-bash:
	sudo docker exec -it dcard_mysql bash
docker-mysql:
	sudo docker exec -it dcard_mysql mysql -D dcard -u root -p
docker-redis:
	sudo docker exec -it dcard_redis redis-cli
migrate-up:
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/dcard" -verbose up
migrate-down:
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/dcard" -verbose down