postgresinit:
	sudo docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:15-alpine

psql:
	sudo docker exec -it postgres15 psql

redis:
	sudo docker exec -it redis7.2 redis-cli

createdb:
	sudo docker exec -it postgres15 createdb --username=root --owner=root go-chat

dropdb:
	sudo docker exec -it postgres15 dropdb go-chat

restart-psql:
	sudo docker restart postgres15

restart-redis:
	sudo docker restart redis7.2

migrateup:
	./migrate -path src/db/migrations -database "postgres://root:root@localhost:5432/go-chat?sslmode=disable" -verbose up

migratedown:
	./migrate -path src/db/migrations -database "postgres://root:root@localhost:5432/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit psql redis createdb dropdb restart-psql restart-redis migrateup migratedown