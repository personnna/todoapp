package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"signinup"
)

func (h *Handler) signUp(c *gin.Context) {
	var input signinup.User

	if err := c.BindJSON(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONP(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"passord" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signinup.User

	if err := c.BindJSON(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONP(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
