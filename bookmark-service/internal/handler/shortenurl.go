package handler

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/repository"
	"github.com/HemlockPham7/golang-system-design/internal/service"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ShortenUrl interface {
	ShortenLink(c *gin.Context)
	Redirect(c *gin.Context)
}

type shortenUrl struct {
	service service.ShortenUrl
}

func NewShortenUrl(svc service.ShortenUrl) ShortenUrl {
	return &shortenUrl{service: svc}
}

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
func (s *shortenUrl) ShortenLink(c *gin.Context) {
	// doc input
	input := &shortenInputBody{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	// goi service de create shorten url
	code, err := s.service.CreateShortenLink(c, input.Url, input.Exp)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenUrl.ShortenLink").Msg("Cannot create shorten url")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})

}

// Redirect Forward the request to the original url
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param code path string true "Shorten code"
// @Success 302
// @Router /v1/links/redirect/{code} [get]
func (s *shortenUrl) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := s.service.GetLinkFromCode(c, code)
	if err != nil {
		if errors.Is(err, repository.ErrCodeNotFound) {
			c.JSON(http.StatusNotFound, response.InputErrResponse)
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenUrl.Redirect").Msg("Cannot get url from code")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.Redirect(http.StatusFound, url)
}
