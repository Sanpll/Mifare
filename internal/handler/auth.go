package handler

import (
	"mifare/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Sign Up
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body dto.SignUpInput true "data to sign up"
// @Success      200  {object}  map[string]interface{} "user ID"
// @Failure      400  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var input dto.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.User.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Sign in (Getting JWT)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body dto.SignInInput true "data to sign in"
// @Success      200  {object}  map[string]interface{} "token"
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var input dto.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
