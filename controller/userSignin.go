package controller

import (
	"net/http"
	"project/initialzer"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SigninInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signin(c *gin.Context) {
	var datas models.User
	var input SigninInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to bind input data",
		})
		return
	}
	if err := initialzer.DB.Where("email=?",input.Email).First(&datas).Error; err != nil{
		c.JSON(http.StatusUnauthorized,gin.H{
			"error":"invalid email",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(datas.Password),[]byte(input.Password)); err != nil{
		c.JSON(http.StatusUnauthorized,gin.H{
			"error":"invalid username or password",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"Successfully logined",
	})
}
