// ══════════════════════════════════════════════
// USER ACCESS HANDLER
// ══════════════════════════════════════════════

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type UserAccessHandler struct {
	accessRepo *repos.UserAccessRepo
}

func NewUserAccessHandler(accessRepo *repos.UserAccessRepo) *UserAccessHandler {
	return &UserAccessHandler{accessRepo: accessRepo}
}



func (h *UserAccessHandler) AssignRole(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.PathValue("id") // 👈 correct Go 1.22+ routing

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	var req struct {
		RoleID int64 `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, "Invalid request body")
		return
	}

	if req.RoleID == 0 {
		utils.ErrorResponse(w, 400, "role_id is required")
		return
	}

	err = h.accessRepo.AssignRoleToUser(userID, req.RoleID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "Role assigned successfully")
}

// GET /api/users/permissions/{userID}
func (h *UserAccessHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/permissions/")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid user ID")
		return
	}

	perms, err := h.accessRepo.GetUserPermissions(userID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{"data": perms})
}

// PUT /api/users/permissions/{userID}
func (h *UserAccessHandler) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/permissions/")
	userID, err := strconv.ParseInt(idStr, 10, 64)
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

	err = h.accessRepo.UpdateUserPermissions(userID, req.Permissions)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "Permissions updated")
}