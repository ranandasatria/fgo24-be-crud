package controllers

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"strconv"

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

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param data body models.User true "User Data"
// @Success 201 {object} utils.Response{results=models.User}
// @Failure 400 {object} utils.Response{success=bool,message=string}
// @Failure 500 {object} utils.Response{success=bool,message=string}
// @Router /user [post]
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
		 fmt.Println("CreateUser error:", err.Error())
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

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := strconv.Atoi(id)

	if err := models.DeleteUser(userID); err != nil {

		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "User ID not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to delete user",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User deleted",
	})
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	userID, _ := strconv.Atoi(id)
	if err := models.UpdateUser(userID, user); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "User ID not found",
			})
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to update user",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User updated",
		Results: user,
	})
}
