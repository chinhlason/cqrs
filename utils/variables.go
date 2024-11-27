package utils

const (
	DBHOST             = "localhost"
	DBPORT             = 5432
	DBUSER             = "root"
	DBPASSWORD         = "root"
	DBNAME             = "cqrs"
	USER_TOPIC         = "debezium.public.users"
	USER_GROUP         = "user-group"
	ORDER_TOPIC        = "debezium.public.orders"
	ORDER_GROUP        = "order-group"
	USER_INDEX         = "users"
	ORDER_INDEX        = "orders"
	KAFKA_BROKER       = "localhost:29092"
	ELASTIC_SEARCH_URL = "http://localhost:9200"
)
