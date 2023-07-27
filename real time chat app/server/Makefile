postgresinit:
	sudo docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:15-alpine

postgres:
	sudo docker exec -it postgres15 psql

createdb:
	sudo docker exec -it postgres15 createdb --username=root --owner=root go-chat

dropdb:
	sudo docker exec -it postgres15 dropdb go-chat

restart:
	sudo docker restart postgres15

migrateup:
	migrate -path src/db/migrations -database "postgres://root:root@localhost:5433/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path src/db/migrations -database "postgres://root:root@localhost:5433/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb restart migrateup migratedown