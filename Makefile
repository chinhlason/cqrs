DSN="postgresql://root:root@localhost:5432/cqrs?sslmode=disable"

g-up:
	goose -dir ./migrations postgres $(DSN) up

g-down:
	goose -dir ./migrations postgres $(DSN) down