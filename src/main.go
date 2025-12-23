package main

import (
	"context"
	"fmt"
	"github.com/alirezamastery/graph_task/db"
	_ "github.com/alirezamastery/graph_task/docs"
	"github.com/alirezamastery/graph_task/middleware"
	"github.com/alirezamastery/graph_task/routes"
	"github.com/alirezamastery/graph_task/utils"
	"log"
	"os"
)

func main() {
	docker := os.Getenv("DOCKER")
	if docker == "" {
		utils.LoadEnvironmentVariables()
	}

	middleware.MustRegisterMetrics()

	shutdown, err := middleware.InitTracing("todo-api")
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	dbConn := db.SetupDB()
	db.InitTasksCount(dbConn)

	db.MigrateDB(dbConn)

	router := routes.SetupRoutes(dbConn)

	apiPort := fmt.Sprintf("0.0.0.0:%s", os.Getenv("API_PORT"))

	err = router.Run(apiPort)
	if err != nil {
		println(err.Error())
		log.Fatalln("error in running router:", err)
	}
}
