package controllers

import (
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func InvalidateUserCache() {
	ctx := context.Background()
	pattern := "*user*"

	iter := utils.RedisClient().Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("Deleting cache key:", iter.Val()) 
		utils.RedisClient().Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Println("Error during Redis SCAN:", err)
	}
}


// GetAllUsers godoc
// @Summary Get all users
// @Tags users
// @Security Token
// @Accept json
// @Produce json
// @Param search query string false "Search by name or email"
// @Param page query string false "Page"
// @Success 200 {object} utils.Response{results=[]models.User}
// @Failure 500 {object} utils.Response{success=bool,message=string}
// @Router /user [get]
func GetAllUsers(ctx *gin.Context) {
	err := utils.RedisClient().Ping(context.Background()).Err()
	noredis := false
	if err != nil && strings.Contains(err.Error(), "refused") {
		noredis = true
	}

	search := ctx.DefaultQuery("search", "")
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("%s?page=%d&search=%s", ctx.FullPath(), page, search)

	if !noredis {
		result := utils.RedisClient().Exists(context.Background(), cacheKey)
		if result.Val() != 0 {
			users := []models.User{}
			data := utils.RedisClient().Get(context.Background(), cacheKey)
			if err = json.Unmarshal([]byte(data.Val()), &users); err == nil {
				ctx.JSON(http.StatusOK, utils.Response{
					Success: true,
					Message: "List users (from Redis)",
					Results: users,
					Page:    page,
					PerPage: 5,
				})
				return
			}
		}
	}

	users, err := models.FindAllUser(search, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	if !noredis {
		encoded, err := json.Marshal(users)
		if err == nil {
			utils.RedisClient().Set(context.Background(), cacheKey, string(encoded), 0)
		}
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "List users (this is from vps)",
		Results: users,
		Page:    page,
		PerPage: 5,
	})
}

// DetailUser godoc
// @Summary Get detail of a user by ID
// @Tags users
// @Security Token
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{results=models.User}
// @Failure 400 {object} utils.Response{success=bool,message=string}
// @Failure 404 {object} utils.Response{success=bool,message=string}
// @Failure 500 {object} utils.Response{success=bool,message=string}
// @Router /user/{id} [get]
func DetailUser(ctx *gin.Context) {
	err := utils.RedisClient().Ping(context.Background()).Err()
	noredis := false
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			noredis = true
		}
	}

	id := ctx.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	if !noredis {
		exists := utils.RedisClient().Exists(context.Background(), ctx.Request.RequestURI)
		if exists.Val() != 0 {
			var cachedUser models.User
			data := utils.RedisClient().Get(context.Background(), ctx.Request.RequestURI)
			if err := json.Unmarshal([]byte(data.Val()), &cachedUser); err == nil {
				ctx.JSON(http.StatusOK, utils.Response{
					Success: true,
					Message: "User detail (from Redis)",
					Results: cachedUser,
				})
				return
			}
		}
	}

	user, err := models.FindUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "User not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to get user data",
			})
		}
		return
	}

	if !noredis {
		encoded, err := json.Marshal(user)
		if err == nil {
			utils.RedisClient().Set(context.Background(), ctx.Request.RequestURI, string(encoded), 0)
		}
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User detail",
		Results: user,
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
	err := utils.RedisClient().Ping(context.Background()).Err()
	noredis := false
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			noredis = true
		}
	}

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

	if !noredis {
	InvalidateUserCache()
}

	ctx.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: "User created",
		Results: user,
	})
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Tags users
// @Security Token
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response{success=bool,message=string}
// @Failure 500 {object} utils.Response{success=bool,message=string}
// @Router /user/{id} [delete]
func DeleteUser(ctx *gin.Context) {
	err := utils.RedisClient().Ping(context.Background()).Err()
	noredis := false
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			noredis = true
		}
	}

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

	if !noredis {
	InvalidateUserCache()
}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User deleted",
	})
}

// UpdateUser godoc
// @Summary Update user data by ID
// @Tags users
// @Security Token
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param data body models.User true "User Data"
// @Success 200 {object} utils.Response{results=models.User}
// @Failure 400 {object} utils.Response{success=bool,message=string}
// @Failure 404 {object} utils.Response{success=bool,message=string}
// @Failure 500 {object} utils.Response{success=bool,message=string}
// @Router /user/{id} [patch]
func UpdateUser(ctx *gin.Context) {
	err := utils.RedisClient().Ping(context.Background()).Err()
	noredis := false
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			noredis = true
		}
	}

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

if !noredis {
	InvalidateUserCache()
}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User updated",
		Results: user,
	})
}
