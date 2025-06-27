package routers

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func userRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.GET("", controllers.GetAllUsers)
	r.GET(":id", controllers.DetailUser)
	r.POST("", controllers.CreateUser)
	r.PATCH(":id", controllers.UpdateUser)
	r.DELETE(":id", controllers.DeleteUser)
}
