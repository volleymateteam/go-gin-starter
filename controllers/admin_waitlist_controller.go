package controllers

import (
	"go-gin-starter/services"
	"go-gin-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAllWaitlist handles GET /admin/waitlist
func GetAllWaitlist(c *gin.Context) {
	waitlist, err := services.GetAllWaitlistEntries()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrDatabase)
		return
	}
	utils.RespondSuccess(c, http.StatusOK, waitlist, utils.MsgWaitlistFetched)
}

// ApproveWaitlistEntry approves a waitlist entry and creates a user
func ApproveWaitlistEntry(c *gin.Context) {
	idParam := c.Param("id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	if err := services.ApproveWaitlistEntry(entryID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgWaitlistApproved)
}

// RejectWaitlistEntry rejects a waitlist entry
func RejectWaitlistEntry(c *gin.Context) {
	idParam := c.Param("id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidUserID)
		return
	}

	if err := services.RejectWaitlistEntry(entryID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, nil, utils.MsgWaitlistRejected)
}
