package link

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/service/link"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Redirect Forward the request to the original url
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param code path string true "Shorten code"
// @Success 302
// @Router /v1/links/redirect/{code} [get]
func (h *linkHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := h.linkService.GetLinkFromCode(c, code)
	if err != nil {
		if errors.Is(err, link.ErrCodeNotFound) {
			c.JSON(http.StatusNotFound, response.InputErrResponse)
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenUrl.Redirect").Msg("Cannot get url from code")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.Redirect(http.StatusFound, url)
}
