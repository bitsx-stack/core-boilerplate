package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)



 


type SettingHandler struct {
	settingRepo *repos.SettingRepo
}

func NewSettingHandler(settingRepo *repos.SettingRepo) *SettingHandler {
	return &SettingHandler{settingRepo: settingRepo}
}

// Get all settings
func (h *SettingHandler) GetAllSettings(response http.ResponseWriter, request *http.Request) {
	
	search := request.URL.Query().Get("search")
	category := request.URL.Query().Get("category")
	limitStr := request.URL.Query().Get("limit")
	offsetStr := request.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ErrorResponse(response, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.ErrorResponse(response, http.StatusBadRequest, err.Error())
		return
	}


	

	result, total, err := h.settingRepo.GetSettings(search, category, limit, offset)
	if err != nil {
		utils.ErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	res := map[string]interface{}{
		"data": result,
		"total": total,
	}

	utils.SuccessResponse(response, http.StatusOK, res)
}

// Get setting by ID
func (h *SettingHandler) GetSetting(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/settings/single/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid setting ID")
		return
	}

	setting, err := h.settingRepo.GetSetting(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if setting == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Setting not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, setting)
}

// CreateSetting - Create a new setting
func (h *SettingHandler) CreateSetting(w http.ResponseWriter, r *http.Request) {
	var req models.Setting
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.Category == "" || req.Code == "" || req.Value == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Code, Name, Value and Category are required")
		return
	}

	err := h.settingRepo.CreateSetting(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Setting created successfully")
}

// UpdateSetting - Update an existing setting
func (h *SettingHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/settings/update/")


	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid setting ID")
		return
	}

	var req models.Setting


	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req.ID = int64(id)
	log.Println(req)

	if   req.Code == "" || req.Name == "" || req.Category == "" || req.Value == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "All fields are required")
		return
	}

	err = h.settingRepo.UpdateSetting(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Setting updated successfully")
}

// DeleteSetting - Soft delete a setting
func (h *SettingHandler) DeleteSetting(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/settings/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid setting ID")
		return
	}

	// For simplicity, using a hardcoded user ID for deletedBy. In a real application, this should come from the authenticated user context.
	deletedBy := int64(1)

	err = h.settingRepo.DeleteSetting(id, deletedBy)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Setting deleted successfully")
}
