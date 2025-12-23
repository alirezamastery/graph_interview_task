package routes

import (
	"github.com/alirezamastery/graph_task/controllers/swagger"
	todoctrl "github.com/alirezamastery/graph_task/controllers/todo"
	"github.com/alirezamastery/graph_task/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"os"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

	middleware.SetupMiddlewares(router)

	// API Routes:
	apiRouter := router.Group("/api")

	todo := todoctrl.NewTodoController(db)
	todoRouter := apiRouter.Group("task")
	{
		todoRouter.GET("/todos/", todo.GetTodoItemList())
		todoRouter.POST("/todos/", todo.CreateTodo())
		todoRouter.GET("/todos/:id", todo.GetTodoItemByID())
		todoRouter.PATCH("/todos/:id", todo.UpdateTodoItem())
		todoRouter.DELETE("/todos/:id", todo.DeleteTodoItem())
	}

	// Swagger:
	swagger.Config()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
