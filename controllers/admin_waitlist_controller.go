package controllers

import (
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitWaitlist handles a new waitlist entry
func SubmitWaitlist(c *gin.Context) {
	var input struct {
		Email  string `json:"email" binding:"required,email"`
		Source string `json:"source"` // Optional
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	if err := services.SubmitWaitlistEntry(input.Email, input.Source); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgWaitlistSuccess)
}

// GetAllWaitlist handles admin viewing of waitlist entries
func GetAllWaitlist(c *gin.Context) {
	emails, err := services.GetAllWaitlistEntries()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{"waitlist_emails": emails}, utils.MsgWaitlistSuccess)
}
