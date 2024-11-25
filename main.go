package main

import (
	"cqrs-postgres-elastic-search-debezium/command"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	db       = "cqrs"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, password, db, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	//command site
	repo := command.NewRepository(db)
	service := command.NewService(repo)
	handler := command.NewHandler(service)

	//router-command-site
	http.HandleFunc("/user/create", handler.InsertUser)
	http.HandleFunc("/user/update", handler.UpdateUser)

	http.HandleFunc("/order/create", handler.InsertOrder)
	http.HandleFunc("/order/update", handler.UpdateOrder)

	log.Println("server started at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
