package user

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/service/user"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login Authentication endpoint
// @Summary Return a jwt token if the input is correct
// @Description Return a jwt token if the input is correct
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param input body loginInput true "Input required"
// @Success 200 {object} object{token=string,message=string} "Success"
// @Router /v1/users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	// doc body input
	input := &loginInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	// call service -> service tra ve token
	token, err := h.service.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, user.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{Message: "invalid credentials"})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to login user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InputErrResponse)
		return
	}

	// tra ve token
	c.JSON(http.StatusOK, response.Message{Message: token})
}
