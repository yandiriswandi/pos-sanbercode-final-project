package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/services"
)

func Authorize(requiredLevels ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenString := c.GetHeader("Authorization")
		// if tokenString == "" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		// 	return
		// }

		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(services.SecretKey), nil
		// })
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix from the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure that the token is using the correct signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(services.SecretKey), nil // Replace with your JWT secret
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Set userID to context
		userID, ok := claims["id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}
		c.Set("userID", int(userID)) // <--- INI YANG PENTING

		// Kalau requiredLevels kosong, skip level check
		if len(requiredLevels) == 0 {
			c.Next()
			return
		}

		// Ambil level user dari token
		userLevel, ok := claims["level"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid level in token"})
			return
		}

		// Cek apakah level user termasuk dalam requiredLevels
		for _, level := range requiredLevels {
			if int(userLevel) == level {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}
