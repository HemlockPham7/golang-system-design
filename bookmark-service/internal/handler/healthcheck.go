package handler

import (
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/service"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HealthCheck interface {
	HealthCheck(c *gin.Context)
}

type healthCheck struct {
	service service.HealthCheck
}

func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheck{service: svc}
}

// HealthCheck checks the health of the service
// @Summary check redis health
// @Description ping and pong with redis server
// @Tags health-check
// @Accept application/json
// @Produce application/json
// @Success 200 {object} model.HealthCheckResponse
// @Router /health-check [get]
func (h *healthCheck) HealthCheck(c *gin.Context) {
	msg, err := h.service.HealthCheck(c)
	if err != nil {
		log.Error().Err(err).Msg("Health-check error")
		c.JSON(http.StatusInternalServerError, response.InstanseErrResponse)
		return
	}
	c.JSON(http.StatusOK, msg)
}
