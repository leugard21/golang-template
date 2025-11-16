package user

import (
	"database/sql"
	"errors"
	"net/http"
	"template/types"
	"template/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if existing, err := h.store.GetUserByEmail(payload.Email); err == nil && existing != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("email already exists"))
		return
	}

	if existing, err := h.store.GetUserByUsername(payload.Username); err == nil && existing != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("username already exists"))
		return
	}

	hashed, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user := types.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashed,
	}

	userID, err := h.store.CreateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message":      "registered successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(payload.Identifier)
	if errors.Is(err, sql.ErrNoRows) {
		u, err = h.store.GetUserByUsername(payload.Identifier)
	}

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	if !utils.CheckPassword(u.Password, payload.Password) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	accessToken, err := utils.GenerateAccessToken(u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":      "login successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
