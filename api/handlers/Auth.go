package handlers

import (
	"encoding/json"
	"net/http"

	"bitpos.com/m/config"
	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type AuthHandler struct {
	userRepo *repos.UserRepo
	config   *config.Config
}

func NewAuthHandler(userRepo *repos.UserRepo, config *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		config:   config,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Credentials")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, h.config.JWTSecret)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, models.LoginResponse{
		Token: token,
	})
}

// func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) bool {

// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		http.Error(w, "Authorization header required", http.StatusUnauthorized)
// 		return false
// 	}

// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || parts[0] != "Bearer" {
// 		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
// 		return false
// 	}

// 	token := parts[1]
// 	_, err := utils.ValidateToken(token, h.config.JWTSecret)
// 	if err != nil {
// 		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
// 		return false
// 	}
// 	return true
// }
