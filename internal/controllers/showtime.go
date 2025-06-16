package controllers

import (
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/responses"
	"bmt_showtime_service/internal/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	var req request.AddShowtimeReq
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

func (s *ShowtimeController) GetShowTime(c *gin.Context) {
	showtimeId, err := strconv.Atoi(c.Param("showtime_id"))
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid showtime id (%s) and request: %v", c.Param("showtime_id"), err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	showtime, status, err := s.ShowtimeService.GetShowtime(ctx, int32(showtimeId))
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "get showtime perform successfully", showtime)
}

func (s *ShowtimeController) GetAllShowTimesByFilmIdInOneDate(c *gin.Context) {
	filmIdStr := c.Query("film_id")
	showDateStr := c.Query("show_date")

	if filmIdStr == "" || showDateStr == "" {
		responses.FailureResponse(c, http.StatusBadRequest, "film_id and show_date are required")
		return
	}

	filmId, err := strconv.Atoi(filmIdStr)
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid film id (%s): %v", filmIdStr, err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	showtimes, status, err := s.ShowtimeService.GetAllShowtimesByFilmIdInOneDate(ctx,
		request.GetAllShowtimesByFilmIdInOneDateReq{
			FilmId:   int32(filmId),
			ShowDate: showDateStr,
		})
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "get all showtimes perform successfully", showtimes)
}

func (s *ShowtimeController) ReleaseShowtime(c *gin.Context) {
	var req request.ReleaseShowtimeByIdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	changedBy := c.GetString(global.X_USER_EMAIL)
	req.ChangedBy = changedBy

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	status, err := s.ShowtimeService.ReleaseShowtime(ctx, req)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "release showtime perform successfully", nil)
}

func (s *ShowtimeController) GetAllFilmsCurrentlyShowing(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	films, status, err := s.ShowtimeService.GetAllFilmsCurrentlyShowing(ctx)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "get all films currently showing successfully", films)
}

func (s *ShowtimeController) GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRange(c *gin.Context) {
	filmIdStr := c.Query("film_id")
	cinemaIdStr := c.Query("cinema_id")

	filmId, err := strconv.Atoi(filmIdStr)
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid film id (%s): %v", filmIdStr, err))

		return
	}

	cinemaId, err := strconv.Atoi(cinemaIdStr)
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid cinema id (%s): %v", cinemaIdStr, err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	showtimes, status, err := s.ShowtimeService.GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRange(ctx,
		request.GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRangeReq{
			FilmId:   int32(filmId),
			CinemaId: int32(cinemaId),
		})
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "get all showtimes successfully", showtimes)
}


