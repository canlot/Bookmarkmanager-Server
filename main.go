package main

import (
	"Bookmarkmanager-Server/Configuration"
	"Bookmarkmanager-Server/Models"
	"Bookmarkmanager-Server/Test"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

var router *gin.Engine

func main() {
	Configuration.Environment = Configuration.Test
	Models.DatabaseConfig()
	setUpTestData()
	gin.SetMode(gin.ReleaseMode)

	router = gin.Default()
	InitializeRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}

func setUpTestData() {
	if Configuration.Environment == Configuration.Test || Configuration.Environment == Configuration.Debug {
		Test.PopulateDatabase()
	}
}
