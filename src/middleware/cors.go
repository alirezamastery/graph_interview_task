package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(engine *gin.Engine) {
	config := cors.DefaultConfig()

	config.AllowCredentials = true
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
	}

	engine.Use(cors.New(config))
}
