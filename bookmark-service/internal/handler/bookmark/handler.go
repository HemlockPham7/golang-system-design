package bookmark

import (
	"github.com/HemlockPham7/golang-system-design/internal/service/bookmark"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateBookmark(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
}

type bookmarkHandler struct {
	bookmarkService bookmark.Service
}

func NewHandler(bookmarkService bookmark.Service) Handler {
	return &bookmarkHandler{bookmarkService: bookmarkService}
}
