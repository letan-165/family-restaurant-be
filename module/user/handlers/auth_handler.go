package handlers

import (
	"context"
	"encoding/json"
	"log"
	"myapp/config/oauth"
	"myapp/module/user/models"
	"myapp/module/user/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	url := oauth.GoogleOauthConfig.AuthCodeURL(oauth.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauth.OauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State không hợp lệ"})
		return
	}

	code := c.Query("code")
	token, err := oauth.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể lấy token"})
		return
	}

	client := oauth.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể lấy thông tin user"})
		return
	}
	defer resp.Body.Close()

	var userInfo models.User
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(500, gin.H{"error": "Decode user info lỗi"})
		return
	}
	userInfo.Role = models.CUSTOMER

	services.CreateUser(userInfo)

	jwtToken, err := services.GenerateToken(userInfo.ID, string(userInfo.Role))
	if err != nil {
		c.JSON(500, gin.H{"error": "Tạo JWT thất bại"})
		return
	}
	log.Println("Generated JWT:", jwtToken)
	c.Header("Authorization", "Bearer "+jwtToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"email":   userInfo.Email,
	})

}
