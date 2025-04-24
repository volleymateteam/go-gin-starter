package controllers

import (
	"go-gin-starter/models"
	"go-gin-starter/services"

	"net/http"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}


func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	user, err := services.CreateUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}


func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "valid-jwt-token"})
}

func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user": "You are authorized!"})
}
