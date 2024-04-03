docker-up:
	sudo docker compose up -d
docker-down:
	sudo docker compose down
docker-bash:
	sudo docker exec -it dcard_mysql bash
migrate-up:
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/dcard" -verbose up