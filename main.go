package main

import (
	"Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	Models.DatabaseConfig(Models.Sqlite, Models.Debug)
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	InitializeRoutes()
	router.Run()
}
