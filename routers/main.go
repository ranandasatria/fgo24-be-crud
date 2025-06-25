package routers

import "github.com/gin-gonic/gin"

func CombineRouter(r *gin.Engine) {
	authRouter(r.Group("/login"))
	userRouter(r.Group("/user"))
	uploadRouter(r.Group("/upload"))
	// tokenRouter(r.Group("/token"))
}
