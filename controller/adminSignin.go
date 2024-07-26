package controller

import (
	"net/http"
	"project/initialzer"
	"project/models"

	"github.com/gin-gonic/gin"
)

type adminInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var input adminInput
	var data models.Admin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		return
	}
	if err := initialzer.DB.Where("email=?", input.Email).First(&data).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email",
		})
		return
	}
	if input.Password == data.Password {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})

}
