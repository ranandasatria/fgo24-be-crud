package routers

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup) {
	r.POST("", controllers.AuthLogin)
}
