package middleware

import "github.com/gin-gonic/gin"

func SetupMiddlewares(engine *gin.Engine) {
	CorsMiddleware(engine)
}
