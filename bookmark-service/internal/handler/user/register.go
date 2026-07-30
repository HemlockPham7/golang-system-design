package user

import (
	"net/http"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type registerInput struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,gte=8"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
}

type userResponse struct {
	Data    *model.User `json:"data"`
	Message string      `json:"message"`
}

// Register handles user registration requests.
// It validates the JSON input, delegates to the service layer for user creation,
// and returns the created user or an appropriate error response.
//
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags User
// @Accept json
// @Produce json
// @Param user body registerInput true "User registration details"
// @Success 201 {object} userResponse
// @Router /v1/users/register [post]
func (h *userHandler) Register(c *gin.Context) {
	input := &registerInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	createdUser, err := h.service.CreateUser(c, input.Username, input.Password, input.DisplayName, input.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InstanseErrResponse)
	}

	c.JSON(http.StatusCreated, &userResponse{
		Data:    createdUser,
		Message: "User created successfully",
	})
}
