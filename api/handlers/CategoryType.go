package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type CategoryTypeHandler struct {
	categoryTypeRepo *repos.CategoryTypeRepo
}

func NewCategoryTypeHandler(categoryTypeRepo *repos.CategoryTypeRepo) *CategoryTypeHandler {
	return &CategoryTypeHandler{categoryTypeRepo: categoryTypeRepo}
}

// PataCategoryTypes - Get all category types
func (h *CategoryTypeHandler) GetCategoryTypes(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit == 0 {
		limit = 10
	}

	data, total, err := h.categoryTypeRepo.GetCategoryTypes(search, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": total,
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

// Get category type by ID
func (h *CategoryTypeHandler) GetCategoryType(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	categoryType, err := h.categoryTypeRepo.GetCategoryTypeByID(int64(id))
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if categoryType == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Category type not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, categoryType)
}



// CreateCategoryType - Create a new category type
func (h *CategoryTypeHandler) CreateCategoryType(w http.ResponseWriter, r *http.Request) {
	var req models.CategoryType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Name == "" || req.Description == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name,and Description are required")
		return
	}

	req.CreatedBy = r.Context().Value("userID").(int64)

	err := h.categoryTypeRepo.CreateCategoryType(req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Category Type Created Successfully!")
}

// UpdateCategoryType - Update an existing category type
func (h *CategoryTypeHandler) UpdateCategoryType(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.CategoryType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Name == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}

	categoryType, err := h.categoryTypeRepo.UpdateCategoryType(int64(id), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := fmt.Sprintf("%s Category Type Updated Successfully", categoryType.Name)

	utils.SuccessResponse(w, http.StatusOK, message)
}

// DeleteCategoryType - Delete a category type
func (h *CategoryTypeHandler) DeleteCategoryType(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid category type ID")
		return
	}

	if err := h.categoryTypeRepo.DeleteCategoryType(int64(id), r.Context().Value("userID").(int64)); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Category type deleted successfully"})
}
