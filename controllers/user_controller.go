package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "valid-jwt-token"})
}

func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user": "You are authorized!"})
}
