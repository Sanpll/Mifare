package handler

import (
	"mifare/internal/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary      Create new transaction
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body dto.CreateTransactionInput true "data to create transaction"
// @Success      200  {object}  map[string]interface{} "transaction ID"
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/transactions [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	var input dto.CreateTransactionInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Transaction.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get all transactions
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.TransactionsResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/transactions [get]
func (h *Handler) GetTransactions(c *gin.Context) {
	transactions, err := h.services.Transaction.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.TransactionsResponse{
		Transactions: transactions,
	})
}

// @Summary      Get transaction by ID
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "transaction ID"
// @Success      200  {object}  dto.TransactionResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/transactions/{id} [get]
func (h *Handler) GetTransactionById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
		return
	}

	transaction, err := h.services.Transaction.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.TransactionResponse{
		ID:                   transaction.ID,
		CardNumber:           transaction.CardNumber,
		Price:                transaction.Price,
		TerminalSerialNumber: transaction.TerminalSerialNumber,
	})
}

// @Summary      Update transaction by ID
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int                     true  "transaction ID"
// @Param        input body   dto.TransactionUpdate  true  "data to update transaction"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      404  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/transactions/{id} [put]
func (h *Handler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var input dto.TransactionUpdate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Transaction.Update(id, input)
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

// @Summary      Delete transaction by ID
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "transaction ID"
// @Success      200  {object}  handler.statusResponse
// @Failure      400  {object}  handler.errorResponse
// @Failure      401  {object}  handler.errorResponse
// @Failure      403  {object}  handler.errorResponse
// @Failure      500  {object}  handler.errorResponse
// @Router       /api/v1/transactions/{id} [delete]
func (h *Handler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
		return
	}

	err = h.services.Transaction.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}
