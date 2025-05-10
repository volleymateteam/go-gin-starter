package controllers

import (
	"go-gin-starter/services"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAllWaitlist handles GET /admin/waitlist
func GetAllWaitlist(c *gin.Context) {
	waitlist, err := services.GetAllWaitlistEntries()
	if err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, constants.ErrDatabase)
		return
	}
	httpPkg.RespondSuccess(c, http.StatusOK, waitlist, constants.MsgWaitlistFetched)
}

// ApproveWaitlistEntry approves a waitlist entry and creates a user
func ApproveWaitlistEntry(c *gin.Context) {
	idParam := c.Param("id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := services.ApproveWaitlistEntry(entryID); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgWaitlistApproved)
}

// RejectWaitlistEntry rejects a waitlist entry
func RejectWaitlistEntry(c *gin.Context) {
	idParam := c.Param("id")
	entryID, err := uuid.Parse(idParam)
	if err != nil {
		httpPkg.RespondError(c, http.StatusBadRequest, constants.ErrInvalidUserID)
		return
	}

	if err := services.RejectWaitlistEntry(entryID); err != nil {
		httpPkg.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(c, http.StatusOK, nil, constants.MsgWaitlistRejected)
}
