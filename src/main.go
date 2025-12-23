package main

import (
	"fmt"
	"github.com/alirezamastery/graph_task/db"
	_ "github.com/alirezamastery/graph_task/docs"
	"github.com/alirezamastery/graph_task/routes"
	"github.com/alirezamastery/graph_task/utils"
	"log"
	"os"
)

// Product Api:
//
//	version: 0.1
//	title: Product Api
//
// Schemes: http, https
// Host:
// BasePath: /api/v1
//
//	Consumes:
//	 - application/json
//
// Produces:
//   - application/json
//
// swagger:meta
func main() {
	docker := os.Getenv("DOCKER")
	if docker == "" {
		utils.LoadEnvironmentVariables()
	}

	dbConn := db.SetupDB()

	db.MigrateDB(dbConn)

	router := routes.SetupRoutes(dbConn)

	apiPort := fmt.Sprintf("0.0.0.0:%s", os.Getenv("API_PORT"))

	err := router.Run(apiPort)
	if err != nil {
		println(err.Error())
		log.Fatalln("error in running router:", err)
	}
}
