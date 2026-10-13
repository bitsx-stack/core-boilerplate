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

type CategoryHandler struct {
	categoryRepo *repos.CategoryRepo
}

func NewCategoryHandler(categoryRepo *repos.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

// GetCategories - Get all categories
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	category_type := r.URL.Query().Get("type")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit == 0 {
		limit = 10
	}

	data, total, err := h.categoryRepo.GetCategories(search, category_type, limit, offset)
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

// PataCategory - Get category by ID
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.categoryRepo.GetCategory(int64(id))
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve category")
		return
	}
	if category == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, category)
}

// CreateCategory - Create a new category
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {

	var req models.Category

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.CategoryTypeID == 0 || req.Description == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name, Description, and Type are required")
		return
	}

	req.CreatedBy = r.Context().Value("userID").(int64)

	err := h.categoryRepo.CreateCategory(req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Category created successfully")
}

// UpdateCategory - Update an existing category
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.CategoryTypeID == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name and category_type_id are required")
		return
	}

	category, err := h.categoryRepo.UpdateCategory(int64(id), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, fmt.Sprintf("%s update successfully", category.Name))
}

// DeleteCategory - Delete a category
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := r.Context().Value("userID").(int64)

	if err := h.categoryRepo.DeleteCategory(int64(id), userID); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
