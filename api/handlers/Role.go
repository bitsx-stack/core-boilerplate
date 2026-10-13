package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type RoleHandler struct {
	roleRepo *repos.RoleRepo
}

func NewRoleHandler (roleRepo *repos.RoleRepo) *RoleHandler {
	return  &RoleHandler{
		roleRepo: roleRepo,
	}
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name        string   `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, "invalid request")
		return
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		utils.ErrorResponse(w, 401, "unauthorized")
		return
	}


	err := h.roleRepo.CreateRole(req.Name, userID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 201, "Role created")
}


func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}

	id := r.PathValue("id")

	var roleID int64
	fmt.Sscan(id, &roleID)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, "invalid request")
		return
	}

	err := h.roleRepo.UpdateRole(req.Name, req.Permissions, roleID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "role updated")
}

func (h *RoleHandler) GetAllRoles(response http.ResponseWriter, request *http.Request) {
	search := request.URL.Query().Get("search")

	limit := 10
	offset := 0

	if v, err := strconv.Atoi(request.URL.Query().Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(request.URL.Query().Get("offset")); err == nil && v >= 0 {
		offset = v
	}

	result, total, err := h.roleRepo.GetAllRoles(search, limit, offset)
	if err != nil {
		utils.ErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(response, http.StatusOK, map[string]interface{}{
		"data":  result,
		"total": total,
	})
}

func (h *RoleHandler) GetSingleRole(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/roles/single/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetSingleRole(id)
	if err != nil {
		utils.ErrorResponse(w, 404, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, role)
}

// ══════════════════════════════════════════════
// ROLE HANDLER — add these methods
// ══════════════════════════════════════════════

func (h *RoleHandler) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/roles/permissions/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid role ID")
		return
	}

	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	err = h.roleRepo.UpdatePermissions(id, req.Permissions)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "Permissions updated")
}

func (h *RoleHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/roles/permissions/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid role ID")
		return
	}

	perms, err := h.roleRepo.GetRolePermissions(id)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{"data": perms})
}