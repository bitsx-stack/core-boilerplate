package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type StoreHandler struct {
	storeRepo *repos.StoreRepo
}

func NewStoreHandler(storeRepo *repos.StoreRepo) *StoreHandler {
	return &StoreHandler{storeRepo: storeRepo}
}

func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var req models.Store
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.Type == "" || req.Category == "" || req.Location == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name, Type, Category and Location are required")
		return
	}

	req.Name = strings.ToUpper(req.Name)
	req.CreatedBy = r.Context().Value("userID").(int64)

	err := h.storeRepo.CreateStore(req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Store Created Successfully!")
}

func (h *StoreHandler) GetAllStores(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	store_type := r.URL.Query().Get("type")
	category := r.URL.Query().Get("category")
	location := r.URL.Query().Get("location")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	stores, total, err := h.storeRepo.GetAllStores(search, store_type, category, location, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  stores,
		"total": total,
	})
}

func (h *StoreHandler) GetSingleStore(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/stores/single/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := h.storeRepo.GetSingleStore(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if store == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Store not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

func (h *StoreHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/stores/update/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid store ID")
		return
	}

	var req models.Store
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.Type == "" || req.Category == "" || req.Location == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name, Type, Category and Location are required")
		return
	}

	req.Name = strings.ToUpper(req.Name)

	err = h.storeRepo.UpdateStore(id, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Store Updated Successfully!")
}

func (h *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/stores/delete/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid store ID")
		return
	}

	deletedBy := r.Context().Value("userID").(int64)

	err = h.storeRepo.DeleteStore(id, deletedBy)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Store Deleted Successfully!")
}