package handlers

import (
	"context"
	"encoding/json"
	"myapp/common/utils"
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

	user, err := services.CreateOrGetUser(userInfo)
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	jwtToken, err := services.GenerateToken(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Tạo JWT thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"email":   userInfo.Email,
		"token":   jwtToken,
	})
}

func InspectToken(c *gin.Context) {
	claims, err := services.InspectToken(c.Param("token"))
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Xác thực",
		"claims":  claims,
	})
}

func GenerateTokenAdmin(c *gin.Context) {
	var request models.GenerateTokenAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSONError(c, err)
		return
	}
	token, err := services.GenerateTokenAdmin(request.Email, request.Secret)
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Xác thực",
		"claims":  token,
	})
}
