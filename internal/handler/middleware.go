package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const(
	authHeader = "Authorization"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, username, isAdmin, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("user_id", userId)
	c.Set("username", username)
	c.Set("is_admin", isAdmin)
	c.Next()
}

func (h *Handler) adminOnly(c *gin.Context) {
    isAdmin, ok := c.Get("is_admin")
    if !ok || !isAdmin.(bool) {
        newErrorResponse(c, http.StatusForbidden, "admin access required")
        c.Abort()
        return
    }
    c.Next()
}

func (h *Handler) isSameUserOrAdmin(c *gin.Context) {
    idURL, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "invalid user_id in URL")
        return
    }

    userID, ok := c.Get("user_id")
    if !ok {
        newErrorResponse(c, http.StatusUnauthorized, "user_id not found in context")
        return
    }

    isAdmin, _ := c.Get("is_admin")
    if !isAdmin.(bool) && userID.(int) != idURL {
		newErrorResponse(c, http.StatusForbidden, "you can only edit your own profile")
		c.Abort()
		return
    }

	c.Next()
}
