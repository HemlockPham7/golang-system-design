package bookmark

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/pkg/requestutils"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// DeleteBookmarkByID generates a Gin framework handler that deletes a bookmark for the authenticated user.
// @Summary      Delete a new bookmark
// @Description  Delete a new bookmark for the authenticated user
// @Tags         Bookmarks
// @Accept       application/json
// @Produce      application/json
// @Param        id       path      string  true  "Bookmark ID"
// @Success      200      {object}  response.Message
// @Failure      400      {object}  response.Message
// @Security     BearerAuth
// @Router       /v1/bookmarks/{id} [delete]
func (h *bookmarkHandler) DeleteBookmarkByID(c *gin.Context) {
	uid, err := requestutils.GetUserIDFromRequest(c)
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

	err = h.bookmarkService.DeleteBookmarkByID(c, uid, bookmarkID)
	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("operation", "DeleteBookmark").Msg("service return error when delete bookmarks")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, response.Message{
		Message: "Success",
	})
}
