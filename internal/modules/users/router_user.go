package user

import (

	"github.com/gin-gonic/gin"
)

func (h *HandlerUser) RoutesUser(rg *gin.RouterGroup){
	users := rg.Group("/user")
	{
		// Register UserRoutes In Here !!
		users.POST("/",h.CreateUser)
	}
}