package controllers

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(ctx *gin.Context) {

	// page, _ := strconv.Atoi(ctx.DefaultQuery(""))

	userIDRaw, _ := ctx.Get("userID")
	userID := int(userIDRaw.(float64))
	fmt.Printf("User yang sedang login adalah %d\n", userID)


	search := ctx.DefaultQuery("search", "")
	users, err := models.FindAllUser(search)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "List users",
		Results: users,
	})
}

func DetailUser(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Detail user",
		Results: map[string]string{
			"id": id,
		},
	})
}

func CreateUser(ctx *gin.Context) {
	fmt.Println("masuk")
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	if err := models.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: "User created",
		Results: user,
	})
}
