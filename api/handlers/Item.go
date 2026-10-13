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

type ItemHandler struct {
	itemRepo *repos.ItemRepo
}

func NewItemHandler(itemRepo *repos.ItemRepo) *ItemHandler {
	return &ItemHandler{
		itemRepo: itemRepo,
	}
}

// GET ALL
func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	itemType := strings.TrimSpace(r.URL.Query().Get("type"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))


	var storeId int
	if storeIdStr := r.URL.Query().Get("store_id"); storeIdStr != "" {
		if l, err := strconv.Atoi(storeIdStr); err == nil && l > 0 {
			storeId = l
		}
	}

	log.Println(storeId)

	limit := 10
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if oStr := r.URL.Query().Get("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o >= 0 {
			offset = o
		}
	}

	data, total, err := h.itemRepo.GetAllItems(search, itemType, category, storeId, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"data":  data,
		"total": total,
	})
}

// CREATE
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
    var req struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category string `json:"category"`
			Uom      string `json:"uom"`
			UnitCost string `json:"unit_cost"`
	}

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
    return
	}

	cost, err := strconv.ParseFloat(req.UnitCost, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid unit cost")
		return
	}

	if req.Name == "" || req.Type == "" || req.Category == "" || req.UnitCost == "" || req.Uom == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "name, type, category, uom required")
		return
	}

	item.CreatedBy = r.Context().Value("userID").(int64)
	item.UnitCost = cost
	item.Name = req.Name
	item.Category = req.Category
	item.Type = req.Type
	item.UOM = req.Uom

	err = h.itemRepo.CreateItem(item)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Item created successfully")
}

// GET SINGLE
func (h *ItemHandler) GetSingleItem(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.itemRepo.GetSingleItem(int64(id))
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if item == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "item not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, item)
}


// UPDATE
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Category string `json:"category"`
		Uom      string `json:"uom"`
		UnitCost json.Number `json:"unit_cost"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.Type == "" || req.Category == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "name, type, category are required")
		return
	}

	cost, err := req.UnitCost.Float64()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Unit Cost")
		return
	}

	item := models.Item{
		Name:     req.Name,
		Type:     req.Type,
		Category: req.Category,
		UOM:      req.Uom,
		UnitCost: cost,
	}

	err = h.itemRepo.UpdateItem(int64(id), item)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Item updated successfully")
}

// DELETE (soft delete)
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	deletedBy := r.Context().Value("userID").(int64)

	err = h.itemRepo.DeleteItem(int64(id), deletedBy)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Item deleted successfully")
}