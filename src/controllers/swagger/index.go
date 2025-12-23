package swagger

import "github.com/alirezamastery/graph_task/docs"

func Config() {
	docs.SwaggerInfo.Title = "Graph Task"
	docs.SwaggerInfo.Description = "Graph Test Task Project"
	docs.SwaggerInfo.BasePath = "/api"
}
