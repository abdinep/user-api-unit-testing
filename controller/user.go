package controller

import (
	"net/http"
	"project/initialzer"
	"project/models"
	"strconv"

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

type EditUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func EditUser(c *gin.Context) {
	var user models.User
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	edit := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}
	if err := initialzer.DB.Model(&models.User{}).Where("id = ?", id).Updates(edit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated user",
		"code":    200,
	})
}
