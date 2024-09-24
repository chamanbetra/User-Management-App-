package main

import (
	"github.com/chamanbetra/user-management-app/config"
	"github.com/chamanbetra/user-management-app/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	database.Connect()

	router := mux.NewRouter()

}
