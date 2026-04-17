package handler

import (
	"mifare/internal/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary      Create new card
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body dto.CreateCardInput true "data to create card"
// @Success      200  {object}  map[string]interface{} "card ID"
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/cards [post]
func (h *Handler) CreateCard(c *gin.Context) {
	var input dto.CreateCardInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	username, ok := c.Get("username")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "username not found in context")
		return
	}

	id, err := h.services.Card.Create(input, username.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get all cards
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.CardsResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/cards [get]
func (h *Handler) GetCards(c *gin.Context) {
	cards, err := h.services.Card.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.CardsResponse{
		Cards: cards,
	})
}

// @Summary      Get card by ID
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "card ID"
// @Success      200  {object}  dto.CardResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/cards/{id} [get]
func (h *Handler) GetCardById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid card id")
		return
	}

	card, err := h.services.Card.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.CardResponse{
		ID:         card.ID,
		CardNumber: card.CardNumber,
		Balance:    card.Balance,
		IsBlocked:  card.IsBlocked,
		OwnerName:  card.OwnerName,
		KeyValue:   card.KeyValue,
	})
}

// @Summary      Update card by ID
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int               true  "card ID"
// @Param        input body   dto.CardUpdate  true  "data to update"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      404  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/cards/{id} [put]
func (h *Handler) UpdateCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid card id")
		return
	}

	var input dto.CardUpdate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Card.Update(id, input)
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

// @Summary      Delete card by ID
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "card ID"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/cards/{id} [delete]
func (h *Handler) DeleteCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid card id")
		return
	}

	err = h.services.Card.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}
