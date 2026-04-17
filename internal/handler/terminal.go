package handler

import (
	"mifare/internal/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary      Create new terminal
// @Tags         Terminals
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body dto.CreateTerminalInput true "data to create terminal"
// @Success      200  {object}  map[string]interface{} "terminal ID"
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminals [post]
func (h *Handler) CreateTerminal(c *gin.Context) {
	var input dto.CreateTerminalInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Terminal.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get all terminals
// @Tags         Terminals
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.TerminalsResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminals [get]
func (h *Handler) GetTerminals(c *gin.Context) {
	terminals, err := h.services.Terminal.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.TerminalsResponse{
		Terminals: terminals,
	})
}

// @Summary      Get terminal by ID
// @Tags         Terminals
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "terminal ID"
// @Success      200  {object}  dto.TerminalResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminals/{id} [get]
func (h *Handler) GetTerminalById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid terminal id")
		return
	}

	terminal, err := h.services.Terminal.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.TerminalResponse{
		ID:           terminal.ID,
		SerialNumber: terminal.SerialNumber,
		Address:      terminal.Address,
		Name:         terminal.Name,
	})
}

// @Summary      Update terminal by ID
// @Tags         Terminals
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int                  true  "terminal ID"
// @Param        input body   dto.TerminalUpdate  true  "data to update terminal"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      404  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminals/{id} [put]
func (h *Handler) UpdateTerminal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid terminal id")
		return
	}

	var input dto.TerminalUpdate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Terminal.Update(id, input)
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

// @Summary      Delete terminal by ID
// @Tags         Terminals
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "terminal ID"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminals/{id} [delete]
func (h *Handler) DeleteTerminal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid terminal id")
		return
	}

	err = h.services.Terminal.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}

// @Summary      Transaction Authorization
// @Tags         Terminal
// @Accept       json
// @Produce      json
// @Param        input body dto.AuthorizeTransactionInput true "auth data"
// @Success      200  {object}  dto.AuthorizeTransactionResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminal/auth [post]
func (h *Handler) AuthorizeTransaction(c *gin.Context) {
	var input dto.AuthorizeTransactionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Transaction.Authorize(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get all keys
// @Tags         Terminal
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.KeysResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/terminal/keys [get]
func (h *Handler) GetAllKeys(c *gin.Context) {
	keys, err := h.services.Key.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.KeysResponse{
		Keys: keys,
	})
}
