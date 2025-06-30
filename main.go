package main

import (
	"backend/routers"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


// @title CRUD Users API
// @version 1.0
// @description Simple user CRUD API with Gin

// @BasePath /

// @SecurityDefinitions.ApiKey Token
// @in header
// @name Authorization


func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Backend is running"})
	})

	routers.CombineRouter(r)
	
	godotenv.Load()

	r.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT")))
}
