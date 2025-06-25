package controllers

import (
	"backend/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateToken(ctx *gin.Context) {
	godotenv.Load()

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": 1,
		"iat":    time.Now().Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Failed to generate token",
		})
	}
	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Token generated",
		Results: token,
	})
}

func VerifyToken(ctx *gin.Context) {

	tokenInput := struct {
		Value string `form:"token"`
	}{}
	ctx.ShouldBind(&tokenInput)

}
