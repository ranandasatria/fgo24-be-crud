package routers

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func userRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetAllUsers)
	r.GET("/:id", controllers.DetailUser)
	r.POST("", controllers.CreateUser)
}
