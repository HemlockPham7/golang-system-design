package bookmark

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/pkg/requestutils"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type updateBookmarkRequest struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// UpdateBookmarkByID generates a Gin framework handler that updates a bookmark for the authenticated user.
// @Summary      Update a new bookmark
// @Description  Update a new bookmark for the authenticated user
// @Tags         Bookmarks
// @Accept       application/json
// @Produce      application/json
// @Param        id       path      string  true  "Bookmark ID"
// @Param        request  body      updateBookmarkRequest  true  "Create Bookmark Request"
// @Success      200      {object}  response.Message
// @Failure      400      {object}  response.Message
// @Security     BearerAuth
// @Router       /v1/bookmarks/{id} [put]
func (h *bookmarkHandler) UpdateBookmarkByID(c *gin.Context) {
	request, uid, err := requestutils.BindInputFromRequestWithAuth[updateBookmarkRequest](c)
	if err != nil {
		return
	}

	bookmarkID := c.Param("id")
	if bookmarkID == "" {
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "Bookmark ID is required",
		})
		return
	}

	_, err = h.bookmarkService.UpdateBookmarkByID(c, request.Description, request.URL, uid, bookmarkID)
	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("operation", "UpdateBookmark").Msg("service return error when update bookmarks")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, response.Message{
		Message: "Success",
	})
}
