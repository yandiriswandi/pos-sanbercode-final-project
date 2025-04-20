package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
	"github.com/yandiriswandi/pos-sanbercode-final-project/services"
)

func Login(c *gin.Context) {
	var loginRequest models.Login

	fmt.Println(loginRequest)
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	user, token, err := services.AuthenticateUser(c, loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"code":    200,
		"user":    user,
		"token":   token,
		"message": `Login success `,
	})
}
