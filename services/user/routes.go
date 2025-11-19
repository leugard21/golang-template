package user

import (
	"database/sql"
	"errors"
	"net/http"
	"template/types"
	"template/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.handleRegister)
	router.POST("/login", h.handleLogin)

}

func (h *Handler) handleRegister(c *gin.Context) {
	var payload types.RegisterUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if existing, err := h.store.GetUserByEmail(payload.Email); err == nil && existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	if existing, err := h.store.GetUserByUsername(payload.Username); err == nil && existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	hashed, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := types.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashed,
	}

	userID, err := h.store.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "registered successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) handleLogin(c *gin.Context) {
	var payload types.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.store.GetUserByEmail(payload.Identifier)
	if errors.Is(err, sql.ErrNoRows) {
		u, err = h.store.GetUserByUsername(payload.Identifier)
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CheckPassword(u.Password, payload.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "login successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
