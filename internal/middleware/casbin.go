package middleware

import (
	"go-starterkit/internal/package/constants"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(constants.CtxUserRole.String())
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid user role",
			})
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method

		allowed, err := enforcer.Enforce(role, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "authorization check failed",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			return
		}

		c.Next()
	}
}
