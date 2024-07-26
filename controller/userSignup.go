package controller

import (
	"log"
	"net/http"
	"project/initialzer"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("failed to bind input datas", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	datas := models.User{Name: input.Name, Email: input.Email, Password: string(hashedPassword)}
	if err := initialzer.DB.Create(&datas).Error; err != nil {
		log.Println("User already exist", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User already exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}
