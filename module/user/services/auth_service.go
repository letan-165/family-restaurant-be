package services

import (
	"errors"
	"fmt"
	"myapp/common/utils"
	"myapp/module/user/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET chưa được cấu hình")
	}

	expireHours := 24
	if val := os.Getenv("JWT_EXPIRE_HOURS"); val != "" {
		if h, err := strconv.Atoi(val); err == nil {
			expireHours = h
		}
	}

	claims := JWTClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func InspectToken(tokenString string) (*JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET chưa được cấu hình")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token không hợp lệ")
}

func GenerateTokenAdmin(email string, secret string) (string, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	user, err := userRepo.FindByEmail(ctx, email)
	if err != nil || user.Email == "" {
		return "", fmt.Errorf("không tìm thấy user với email: %s", email)
	}

	if user.Role != models.ADMIN {
		return "", fmt.Errorf("không có quyền admin")
	}

	claims := JWTClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
