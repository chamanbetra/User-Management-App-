package main

import (
	"log"
	"net/http"

	"github.com/chamanbetra/user-management-app/config"
	"github.com/chamanbetra/user-management-app/database"
	"github.com/chamanbetra/user-management-app/routes"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadEnv()

	database.Connect()

	r := routes.Router()

	log.Println("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
