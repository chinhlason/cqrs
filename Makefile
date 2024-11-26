DSN ?= "postgresql://root:root@localhost:5432/cqrs?sslmode=disable"

TOPIC ?= debezium.public.users
HOST ?= docker-kafka-1

TABLE ?= users


g-up:
	goose -dir ./migrations postgres $(DSN) up

g-down:
	goose -dir ./migrations postgres $(DSN) down

consume:
	docker exec $(HOST) kafka-console-consumer --bootstrap-server localhost:9092 --topic $(TOPIC)

identity:
	docker exec -it postgres psql -U root -d cqrs -c "ALTER TABLE $(TABLE) REPLICA IDENTITY FULL;"

cfg:
	goose -dir ./migrations postgres $(DSN) up
	docker exec -it postgres psql -U root -d cqrs -c "ALTER SYSTEM SET wal_level = 'logical';"
	docker exec -it postgres psql -U root -d cqrs -c "ALTER SYSTEM SET max_replication_slots = 4;"
	docker exec -it postgres psql -U root -d cqrs -c "ALTER SYSTEM SET max_wal_senders = 4;"
	docker restart postgres
	docker exec -it postgres psql -U root -d cqrs -c "ALTER USER root WITH REPLICATION;"
	docker exec -it postgres psql -U root -d cqrs -c "CREATE PUBLICATION cqrs_pub FOR ALL TABLES;"
	docker exec -it postgres psql -U root -d cqrs -c "ALTER TABLE users REPLICA IDENTITY FULL;"
	docker exec -it postgres psql -U root -d cqrs -c "ALTER TABLE orders REPLICA IDENTITY FULL;"
	curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d "@debezium.json"
conn:
	curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d "@debezium.json"

delete:
	curl -i -X DELETE -H "Accept:application/json" localhost:8083/connectors/postgres-connector

up:
	docker-compose up -d

down:
	docker-compose down

run:
	go run main.go

es-up:
	docker-compose -f elastic-search.docker-compose.yml up -d