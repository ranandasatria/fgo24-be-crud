package controllers

import (
	"backend/models"
	"backend/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthLogin(ctx *gin.Context) {
	godotenv.Load()
	secretKey := os.Getenv("APP_SECRET")

	form := struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}{}

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	user, err := models.FindOneUserByEmail(form.Email)

	if err != nil {
		//handle
	}

	if user == (models.User{}) || (form.Password != user.Password) {
		ctx.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Wrong email or password",
		})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"iat": time.Now().Unix(),
	})

	token, _ := generateToken.SignedString([]byte(secretKey))

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Token generated",
		Results: map[string]string{
			"token": token,
		},
	})
}
