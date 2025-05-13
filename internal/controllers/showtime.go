package controllers

import (
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/responses"
	"bmt_showtime_service/internal/services"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ShowtimeController struct {
	ShowtimeService services.IShowtime
}

func NewShowtimeController(
	showtimeService services.IShowtime) *ShowtimeController {
	return &ShowtimeController{
		ShowtimeService: showtimeService,
	}
}

func (s *ShowtimeController) AddShowtime(c *gin.Context) {
	var req request.AddShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	changedBy := c.GetString(global.X_USER_EMAIL)
	req.ChangedBy = changedBy

	status, err := s.ShowtimeService.AddShowtime(ctx, req)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "add new showtime perform successfully", nil)

}
