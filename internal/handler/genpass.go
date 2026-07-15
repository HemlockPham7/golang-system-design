package handler

import (
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/service"
	"github.com/gin-gonic/gin"
)

const passwordLength = 12

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassService service.GenPass
}

func NewGenPass(genPassService service.GenPass) GenPass {
	return &genPassHandler{genPassService: genPassService}
}

func (g *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := g.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Err",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"password": pass,
	})
}
