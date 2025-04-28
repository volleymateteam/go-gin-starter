package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitWaitlist handles POST /api/waitlist/submit
func SubmitWaitlist(c *gin.Context) {
	var input dto.CreateWaitlistEntryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidInput)
		return
	}

	err := services.SubmitWaitlistEntry(input.Email, input.Source)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, nil, utils.MsgWaitlistSubmitted)
}

// GetAllWaitlist handles GET /api/admin/waitlist
func GetAllWaitlist(c *gin.Context) {
	entries, err := services.GetAllWaitlistEntries()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, entries, utils.MsgWaitlistFetched)
}
