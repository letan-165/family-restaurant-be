package middlewares

import (
	"strings"

	"myapp/module/user/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Chưa có token"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Token không hợp lệ"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := services.InspectToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Token không hợp lệ", "detail": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
