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

type CinemaController struct {
	CinemaService services.ICinema
}

func NewCinemaController(
	cinemaService services.ICinema,
) *CinemaController {
	return &CinemaController{
		CinemaService: cinemaService,
	}
}

func (c *CinemaController) GetCinemasForShowingFilm(gc *gin.Context) {
	filmId, err := strconv.Atoi(gc.Param("film_id"))
	if err != nil {
		responses.FailureResponse(gc, http.StatusBadRequest, fmt.Sprintf("invalid film_id (%s)", gc.Param("film_id")))
		return
	}

	ctx, cancel := context.WithTimeout(gc.Request.Context(), 10*time.Second)
	defer cancel()

	cinema, status, err := c.CinemaService.GetCinemasForShowingFilm(ctx, int32(filmId))
	if err != nil {
		responses.FailureResponse(gc, status, err.Error())
		return
	}

	responses.SuccessResponse(gc, status, "get cinema for showing film successfully", cinema)
}
