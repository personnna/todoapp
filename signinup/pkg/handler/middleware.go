package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	autharizationHeader = "Autharization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(autharizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, "")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("userId", userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user if not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type found")
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}