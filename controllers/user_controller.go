package controllers

import (
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"log"

	"net/http"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID		uint   `json:"id"`
	Username string `json:"username"`
	Email 	string `json:"email"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	// Debug log to confirm password check execution
	log.Printf("Checking password strength for: %s", input.Password)
	log.Printf("Password validation result: %v for password: %s", utils.IsStrongPassword(input.Password), input.Password)
	if !utils.IsStrongPassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters."})
		return
	}
	user, err := services.CreateUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reponse := UserResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
	}
	c.JSON(http.StatusCreated, gin.H{"user": reponse})
}



func GetProfile(c *gin.Context) {
	userINRow, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userINRow.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,})
}


func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
