package link

import (
	"net/http"

	"github.com/HemlockPham7/golang-system-design/pkg/requestutils"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type shortenInputBody struct {
	Url string `json:"url" binding:"required,url"`
	Exp int64  `json:"exp" binding:"required,gte=60"`
}

// ShortenLink Generate shorten link
// @Summary Generate shorten url based on original url that last upto 7 days
// @Description Generate shorten url based on original url that last upto 7 days
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param input body shortenInputBody true "Input required"
// @Success 200 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *linkHandler) ShortenLink(c *gin.Context) {
	// doc input
	input, err := requestutils.BindInputFromRequest[shortenInputBody](c)
	if err != nil {
		return
	}

	// goi service de create shorten url
	code, err := h.linkService.CreateShortenLink(c, input.Url, input.Exp)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenUrl.ShortenLink").Msg("Cannot create shorten url")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})

}
