package bookmark_Server

import "github.com/gin-gonic/gin"

var router *gin.Engine

func main() {
	gin.SetMode(gin.ReleaseMode)

	router = gin.Default()

	initializeRoutes()

	router.Run()

}
