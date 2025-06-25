package routers

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func uploadRouter(r *gin.RouterGroup) {
	r.POST("", controllers.UploadFile)
}
