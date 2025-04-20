package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
	"github.com/yandiriswandi/pos-sanbercode-final-project/utils"
)

var SecretKey = []byte{18, 12, 1987 % 256, 1998 % 256} // Menggunakan % 256 untuk memastikan dalam rentang 0-255

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"level":    user.Level,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func AuthenticateUser(c *gin.Context, username, password string) (models.User, string, error) {
	var user models.User

	query := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	err := config.DB.Get(&user, query, username)
	if err != nil {
		return user, "", fmt.Errorf("user not found")
	}
	hashedPassword := utils.HashPassword(password)
	if user.Password != hashedPassword {
		return user, "", fmt.Errorf("invalid credentials")
	}

	token, err := GenerateToken(user)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}
