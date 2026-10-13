package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bitpos.com/m/config"
	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type UserHandler struct {
	userRepo *repos.UserRepo
}

func NewUserHandler(userRepo *repos.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// // PataTaarifa - Get current user profile
// func (h *UserHandler) PataTaarifa(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value("userID").(int64)

// 	user, err := h.userRepo.PataUserKwaID(userID)
// 	if err != nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Mtumiaji hakuweza kupatikana")
// 		return
// 	}
// 	if user == nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusNotFound, "Mtumiaji hayupo katika takwimu zetu")
// 		return
// 	}

// 	rudishaMajibuKwaMpangilio(w, http.StatusOK, user)
// }

// // SahihishaWasifu - Update current user profile
// // func (h *UserHandler) SahihishaWasifu(w http.ResponseWriter, r *http.Request) {
// // 	userID := r.Context().Value("userID").(int64)

// // 	var req models.UpdateProfileRequest
// // 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusBadRequest, "Ombi batili")
// // 		return
// // 	}

// // 	user, err := h.userRepo.PataUserKwaID(userID)
// // 	if err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Mtumiaji hakuweza kupatikana")
// // 		return
// // 	}
// // 	if user == nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusNotFound, "Mtumiaji hayupo katika takwimu zetu")
// // 		return
// // 	}

// // 	user.Name = req.Name
// // 	if err := h.userRepo.SahihishaTaarifaZaUser(user); err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Imeshindwa kusasisha taarifa")
// // 		return
// // 	}

// // 	rudishaMajibuKwaMpangilio(w, http.StatusOK, user)
// // }

// Get user by ID
func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	token := parts[1]
	claims, err := utils.ValidateToken(token, cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, user)
}

// // PataWatumiaji - Get all users (Admin only)
// func (h *UserHandler) PataWatumiaji(w http.ResponseWriter, r *http.Request) {
// 	users, err := h.userRepo.PataWatumiaji()
// 	if err != nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Failed to retrieve users")
// 		return
// 	}

// 	if users == nil {
// 		users = []models.User{}
// 	}

// 	rudishaMajibuKwaMpangilio(w, http.StatusOK, users)
// }

// CreateUser - Create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req *models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Email == "" || req.Name == "" || req.Phone == "" || req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Email, Phone, Password and Name are required")
		return
	}

	existingUser, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existingUser != nil {
		utils.ErrorResponse(w, http.StatusConflict, "Email already exists")
		return
	}

	user := &models.User{
		Email:     req.Email,
		Name:      req.Name,
		Phone:     req.Phone,
		Status:    req.Status,
		CreatedAt: req.CreatedAt,
	}

	// If password is provided, hash it
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to process password")
			return
		}
		user.Password = hashedPassword
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = utils.SendOTP(user.Phone)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "User created successfully!")
}

// // UpdateUser - Update a user (Admin only)
// // func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
// // 	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
// // 	id, err := strconv.ParseInt(idStr, 10, 64)
// // 	if err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusBadRequest, "Invalid user ID")
// // 		return
// // 	}

// // 	var req models.UpdateProfileRequest
// // 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusBadRequest, "Invalid request")
// // 		return
// // 	}

// // 	user, err := h.userRepo.PataMtumiaji(id)
// // 	if err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Failed to retrieve user")
// // 		return
// // 	}
// // 	if user == nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusNotFound, "User not found")
// // 		return
// // 	}

// // 	user.Name = req.Name
// // 	if err := h.userRepo.SahihishaTaarifaZaUser(user); err != nil {
// // 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Failed to update user")
// // 		return
// // 	}

// // 	rudishaMajibuKwaMpangilio(w, http.StatusOK, user)
// // }

// // DeleteUser - Delete a user (Admin only)
// func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusBadRequest, "Invalid user ID")
// 		return
// 	}

// 	user, err := h.userRepo.PataMtumiaji(id)
// 	if err != nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Failed to retrieve user")
// 		return
// 	}
// 	if user == nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusNotFound, "User not found")
// 		return
// 	}

// 	if err := h.userRepo.FutaUser(id); err != nil {
// 		rudishaMajibuKwaMpangilioNaMakosa(w, http.StatusInternalServerError, "Failed to delete user")
// 		return
// 	}

// 	rudishaMajibuKwaMpangilio(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
// }







func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")

	limit := 10
	offset := 0
	if v, err := strconv.Atoi(q.Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(q.Get("offset")); err == nil && v >= 0 {
		offset = v
	}

	data, total, err := h.userRepo.GetAllUsers(search, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{
		"data":  data,
		"total": total,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/single/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}
	if user == nil {
		utils.ErrorResponse(w, 404, "User not found")
		return
	}

	utils.SuccessResponse(w, 200, user)
}

// func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 		Name     string `json:"name"`
// 		Phone    string `json:"phone"`
// 		Status   bool   `json:"status"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.ErrorResponse(w, 400, err.Error())
// 		return
// 	}

// 	if req.Name == "" || req.Phone == "" || req.Email == "" {
// 		utils.ErrorResponse(w, 400, "Name, email and phone are required")
// 		return
// 	}
// 	if req.Password == "" || len(req.Password) < 6 {
// 		utils.ErrorResponse(w, 400, "Password must be at least 6 characters")
// 		return
// 	}

// 	// Hash password
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		utils.ErrorResponse(w, 500, "Failed to hash password")
// 		return
// 	}

// 	err = h.userRepo.CreateUser(models.User{
// 		Email:    req.Email,
// 		Password: string(hashed),
// 		Name:     strings.ToUpper(req.Name),
// 		Phone:    req.Phone,
// 		Status:   req.Status,
// 	})

// 	if err != nil {
// 		utils.ErrorResponse(w, 400, err.Error())
// 		return
// 	}

// 	utils.SuccessResponse(w, 201, "User created successfully")
// }

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/update/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	var req struct {
		Email  string `json:"email"`
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Status bool   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	if req.Name == "" || req.Phone == "" || req.Email == "" {
		utils.ErrorResponse(w, 400, "Name, email and phone are required")
		return
	}

	err = h.userRepo.UpdateUser(id, models.User{
		Email:  req.Email,
		Name:   strings.ToUpper(req.Name),
		Phone:  req.Phone,
		Status: req.Status,
	})

	if err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "User updated successfully")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/delete/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	deletedBy := r.Context().Value("userID").(int64)

	err = h.userRepo.DeleteUser(id, deletedBy)
	if err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "User deleted successfully")
}

// // Assign role to user
// r.HandleFunc("/api/users/assign-role/", h.AssignRole).Methods("PUT")

// func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
//     idStr := strings.TrimPrefix(r.URL.Path, "/api/users/assign-role/")
//     id, _ := strconv.ParseInt(idStr, 10, 64)

//     var req struct {
//         RoleID int64 `json:"role_id"`
//     }
//     json.NewDecoder(r.Body).Decode(&req)

//     _, err := h.db.Exec("UPDATE users SET role_id = $1 WHERE id = $2", req.RoleID, id)
//     if err != nil {
//         utils.ErrorResponse(w, 500, err.Error())
//         return
//     }
//     utils.SuccessResponse(w, 200, "Role assigned")
// }

// // Get user permissions
// r.HandleFunc("/api/users/permissions/", h.GetPermissions).Methods("GET")

// // Save user permissions
// r.HandleFunc("/api/users/permissions/", h.SavePermissions).Methods("PUT")

// // List all available permissions
// r.HandleFunc("/api/permissions", h.ListPermissions).Methods("GET")


// ══════════════════════════════════════════════
// USER HANDLER — add these methods
// ══════════════════════════════════════════════

func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/assign-role/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	var req struct {
		RoleID int64 `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	if req.RoleID == 0 {
		utils.ErrorResponse(w, 400, "role_id is required")
		return
	}

	err = h.userRepo.AssignRole(id, req.RoleID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "Role assigned")
}

func (h *UserHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/permissions/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	perms, err := h.userRepo.GetUserPermissions(id)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{"data": perms})
}

func (h *UserHandler) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/permissions/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	err = h.userRepo.UpdateUserPermissions(id, req.Permissions)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "Permissions updated")
}