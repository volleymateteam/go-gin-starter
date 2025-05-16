package controllers

import (
	"go-gin-starter/dto"
	"go-gin-starter/pkg/constants"
	httpPkg "go-gin-starter/pkg/http"
	"go-gin-starter/pkg/logger"
	"go-gin-starter/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// WaitlistController handles waitlist-related HTTP requests
type WaitlistController struct {
	waitlistService services.WaitlistService
}

// NewWaitlistController creates a new instance of WaitlistController
func NewWaitlistController(waitlistService services.WaitlistService) *WaitlistController {
	return &WaitlistController{
		waitlistService: waitlistService,
	}
}

// SubmitWaitlist handles waitlist submission
func (c *WaitlistController) SubmitWaitlist(ctx *gin.Context) {
	var input dto.CreateWaitlistEntryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := c.waitlistService.SubmitWaitlistEntry(input.Email, input.Source)
	if err != nil {
		if err.Error() == constants.ErrAlreadyInWaitlist {
			httpPkg.RespondError(ctx, http.StatusConflict, constants.ErrAlreadyInWaitlist)
			return
		}
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgWaitlistSuccess)
}

// GetAllWaitlist returns all waitlist entries
func (c *WaitlistController) GetAllWaitlist(ctx *gin.Context) {
	entries, err := c.waitlistService.GetAllWaitlistEntries()
	if err != nil {
		logger.Error("Failed to get waitlist entries", zap.Error(err))
		httpPkg.RespondError(ctx, http.StatusInternalServerError, constants.ErrInternalServer)
		return
	}

	// Always return an initialized array, even if empty
	if entries == nil {
		entries = []dto.WaitlistEntryResponse{}
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, entries, constants.MsgWaitlistFetched)
}

// ApproveWaitlistEntry approves a waitlist entry
func (c *WaitlistController) ApproveWaitlistEntry(ctx *gin.Context) {
	id := ctx.Param("id")
	entryID, err := uuid.Parse(id)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidID)
		return
	}

	if err := c.waitlistService.ApproveWaitlistEntry(entryID); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgWaitlistApproved)
}

// RejectWaitlistEntry rejects a waitlist entry
func (c *WaitlistController) RejectWaitlistEntry(ctx *gin.Context) {
	id := ctx.Param("id")
	entryID, err := uuid.Parse(id)
	if err != nil {
		httpPkg.RespondError(ctx, http.StatusBadRequest, constants.ErrInvalidID)
		return
	}

	if err := c.waitlistService.RejectWaitlistEntry(entryID); err != nil {
		httpPkg.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httpPkg.RespondSuccess(ctx, http.StatusOK, nil, constants.MsgWaitlistRejected)
}
