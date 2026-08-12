package link

import (
	"github.com/HemlockPham7/golang-system-design/internal/service/link"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	ShortenLink(c *gin.Context)
	Redirect(c *gin.Context)
}

type linkHandler struct {
	linkService link.Service
}

func NewLinkHandler(linkService link.Service) Handler {
	return &linkHandler{linkService: linkService}
}
