package controller

import (
	"net/http"
	"project/initialzer"
	"project/models"

	"github.com/gin-gonic/gin"
)

func ListUser(c *gin.Context) {
	var list []models.User
	var listUser []gin.H
	count := 0
	if err := initialzer.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}
	for _, user := range list {
		listUser = append(listUser, gin.H{
			"user_id":    user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		})
		count++
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       listUser,
		"total_user": count,
		"status":     200,
	})
}

type editUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func EditUser(c *gin.Context) {
	var edit editUser
	var user models.User
	userID := c.Param("ID")
	if err := initialzer.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "user not found",
			"status": 401,
		})
		return
	}
	if err := c.ShouldBindJSON(&edit); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to bind data",
		})
		return
	}
	user.Name = edit.Name
	user.Email = edit.Email
	if err := initialzer.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated",
		"status":  200,
	})
}
