package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pubestpubest/g12-todo-backend/database"
	"github.com/pubestpubest/g12-todo-backend/middlewares"
	"github.com/pubestpubest/g12-todo-backend/routes"
	log "github.com/sirupsen/logrus"
)

func init() {
	fmt.Println("Hello, World from init()")

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	runEnv := os.Getenv("RUN_ENV")
	if runEnv == "" {
		runEnv = "development"
	}

	deployEnv := os.Getenv("DEPLOY_ENV")
	if deployEnv == "" {
		deployEnv = "local"
	}

	if deployEnv == "local" {
		if err := godotenv.Load("configs/.env"); err != nil {
			log.Fatal("[init]: Error loading .env file: ", err)
		}
	}

	if runEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info("[init]: Run environment: ", runEnv)

	if err := database.ConnectDB(runEnv); err != nil {
		log.Fatal("[init]: Connect database PG error: ", err.Error())
	}
}

func main() {
	fmt.Println("Hello, World from main()")

	app := gin.Default()

	app.Use(middlewares.CORSMiddleware())

	app.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "not found",
		})
	})
	app.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status": "method not allowed",
		})
	})

	v1 := app.Group("/v1")
	routes.EventRoutes(v1)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}
	app.Run(":" + port)
}
