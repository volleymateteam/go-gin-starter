package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitWaitlist handles POST /api/waitlist/submit
func SubmitWaitlist(c *gin.Context) {
	var input dto.CreateWaitlistEntryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidInput)
		return
	}

	err := services.SubmitWaitlistEntry(input.Email, input.Source)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusCreated, nil, constants.MsgWaitlistSubmitted)
}
