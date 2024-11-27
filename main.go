package main

import (
	"cqrs-postgres-elastic-search-debezium/command"
	"cqrs-postgres-elastic-search-debezium/query"
	"cqrs-postgres-elastic-search-debezium/sync"
	"cqrs-postgres-elastic-search-debezium/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		utils.DBHOST, utils.DBUSER, utils.DBPASSWORD, utils.DBNAME, utils.DBPORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	//command site
	repo := command.NewRepository(db)
	commandService := command.NewService(repo)
	commandHandler := command.NewHandler(commandService)

	//router-command-site
	http.HandleFunc("/user/create", commandHandler.InsertUser)
	http.HandleFunc("/user/update", commandHandler.UpdateUser)

	http.HandleFunc("/order/create", commandHandler.InsertOrder)
	http.HandleFunc("/order/update", commandHandler.UpdateOrder)

	es, err := query.GetESClient(utils.ELASTIC_SEARCH_URL)
	if err != nil {
		log.Fatalf("failed to get es client: %v", err)
	}

	esClient := query.NewESClient(es)

	consume := sync.NewConsumerGroupHandler(esClient)

	go consume.ConsumeMessage(utils.KAFKA_BROKER, utils.USER_TOPIC, utils.USER_GROUP)
	go consume.ConsumeMessage(utils.KAFKA_BROKER, utils.ORDER_TOPIC, utils.ORDER_GROUP)

	queryService := query.NewQueryService(esClient)
	queryHandler := query.NewQueryHandler(queryService)

	//router-query-site
	http.HandleFunc("/search", queryHandler.Search)

	log.Println("server started at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
