package handler

import (
	"net/http"

	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const passwordLength = 12

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassService utils.GenPass
}

func NewGenPass(genPassService utils.GenPass) GenPass {
	return &genPassHandler{genPassService: genPassService}
}

// GeneratePassword godoc
// @Summary Generate a random password
// @Tags password
// @Produce json
// @Success 200 {string} string
// @Router /genpass [get]
func (g *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := g.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.genPassHandler.GeneratePassword").Msg("Cannot generate password")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"password": pass,
	})
}
