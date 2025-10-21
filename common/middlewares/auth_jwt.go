package middlewares

import (
	"net/http"
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
			c.JSON(401, gin.H{"message": "Token không hợp lệ", "error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy role"})
			c.Abort()
			return
		}

		userRole := roleValue.(string)

		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền truy cập"})
		c.Abort()
	}
}
