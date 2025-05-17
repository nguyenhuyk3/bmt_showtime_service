package controllers

import (
	"bmt_showtime_service/internal/responses"
	"bmt_showtime_service/internal/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ShowtimeSeatController struct {
	ShowtimeSeatService services.IShowtimeSeat
}

func NewShowtimeSeatController(
	showtimeSeatService services.IShowtimeSeat) *ShowtimeSeatController {
	return &ShowtimeSeatController{
		ShowtimeSeatService: showtimeSeatService,
	}
}

func (s *ShowtimeSeatController) GetAllShowtimeSeatsByShowtimeId(c *gin.Context) {
	showtimeId, err := strconv.Atoi(c.Query("showtime_id"))
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid showtime id (%s): %v", c.Query("showtime_id"), err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	seats, status, err := s.ShowtimeSeatService.GetAllShowtimeSeatsByShowtimeId(ctx, int32(showtimeId))
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "get all showtime seats perform successfully", seats)
}
