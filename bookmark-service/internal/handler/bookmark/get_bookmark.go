package bookmark

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/dto"
	"github.com/HemlockPham7/golang-system-design/pkg/requestutils"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type getBookmarksRequest struct {
	Page  int `form:"page" validate:"gte=1"`
	Limit int `form:"limit" validate:"gte=1"`
}

// GetBookmarks handles the HTTP request to list bookmarks for the authenticated user.
// It supports pagination and sorting based on query parameters.
// @Summary      List bookmarks
// @Description  Retrieve a list of bookmarks for the authenticated user with pagination and sorting
// @Tags         Bookmarks
// @Accept       application/json
// @Produce      application/json
// @Param        page   query     int    false  "Page number"        default(1)    example(1)
// @Param        limit  query     int    false  "Number of items per page"  default(5)    example(5)
// @Success      200    {object}  object{data=[]model.Bookmark} "Success"
// @Failure      400    {object}  response.Message
// @Failure      401    {object}  response.Message
// @Failure      500    {object}  response.Message
// @Security	 BearerAuth
// @Router       /v1/bookmarks [get]
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	request, uid, err := requestutils.BindInputFromRequestWithAuth[getBookmarksRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	result, err := h.bookmarkService.GetBookmarks(c, uid, request.Page, request.Limit)
	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("operation", "GetBookmarks").Msg("service return error when get bookmarks")
		c.JSON(http.StatusInternalServerError, response.InstanseErrResponse)
		return
	}

	c.JSON(http.StatusOK, &dto.SuccessResponse[[]*model.Bookmark]{
		Data: result.Bookmarks,
		Pagination: &dto.Pagination{
			Page:  request.Page,
			Limit: request.Limit,
			Total: result.Total,
		},
	})
}
