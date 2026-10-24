package middleware

import (
	"go-starterkit/internal/package/constants"
	"go-starterkit/internal/package/jwt"
	"go-starterkit/internal/response"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth memvalidasi Access Token dari Header Authorization
func Auth(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Header otorisasi tidak ditemukan", nil)
			c.Abort()
			return
		}

		// 2. Format harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Format token salah (harus Bearer <token>)", nil)
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// 3. Validasi token via JWTService
		claims, err := jwtService.ValidateAccessToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Access token tidak valid atau kedaluwarsa", err.Error())
			c.Abort()
			return
		}

		// 4. Inject UserID dan Role ke Gin Context
		c.Set(constants.CtxUserID.String(), claims.UserID)
		c.Set(constants.CtxUserRole.String(), claims.Role)

		c.Next()
	}
}
