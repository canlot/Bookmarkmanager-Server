package main

import (
	"github.com/canlot/Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	Models.DatabaseConfig()

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	initializeRoutes()
	router.Run()

}
