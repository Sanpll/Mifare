package handler

import (
	"mifare/internal/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary      Get all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.UsersResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.services.User.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.UsersResponse{
		Users: users,
	})
}

// @Summary      Get user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "user ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/users/{id} [get]
func (h *Handler) GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.services.User.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	})
}

// @Summary      Update user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int            true  "user ID"
// @Param        input body   dto.UserUpdate  true  "data to update user"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      404  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var input dto.UserUpdate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.User.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}

// @Summary      Delete user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "user ID"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.services.User.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}