package controllers

import (
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Upload failed",
		})
		return
	}

	if file.Size > 2*1024*1024 {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "File is too large",
		})
		return
	}


	ctx.SaveUploadedFile(file, "./uploads/"+file.Filename)

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Upload success",
	})
}



	// fileName := uuid.New().String()