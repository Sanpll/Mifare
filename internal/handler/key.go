package handler

import (
	"mifare/internal/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary      Create new key
// @Tags         Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body dto.CreateKeyInput true "data to create key"
// @Success      200  {object}  map[string]interface{} "key ID"
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/keys [post]
func (h *Handler) CreateKey(c *gin.Context) {
	var input dto.CreateKeyInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Key.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get all keys
// @Tags         Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.KeysResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/keys [get]
func (h *Handler) GetKeys(c *gin.Context) {
	keys, err := h.services.Key.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.KeysResponse{
		Keys: keys,
	})
}

// @Summary      Get key by ID
// @Tags         Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "key ID"
// @Success      200  {object}  dto.KeyResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/keys/{id} [get]
func (h *Handler) GetKeyById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid key id")
		return
	}

	key, err := h.services.Key.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.KeyResponse{
		ID:          key.ID,
		KeyValue:    key.KeyValue,
		KeyType:     key.KeyType,
		Description: key.Description,
	})
}

// @Summary      Update key by ID
// @Tags         Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int             true  "key ID"
// @Param        input body   dto.KeyUpdate  true  "data to update key"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      404  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/keys/{id} [put]
func (h *Handler) UpdateKey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid key id")
		return
	}

	var input dto.KeyUpdate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Key.Update(id, input)
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

// @Summary      Delete key by ID
// @Tags         Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "key ID"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/keys/{id} [delete]
func (h *Handler) DeleteKey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid key id")
		return
	}

	err = h.services.Key.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}
