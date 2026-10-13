package handlers

import (
	"net/http"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type BulkHandler struct {
	repo *repos.BulkRepo
}

func NewBulkHandler(r *repos.BulkRepo) *BulkHandler {
	return &BulkHandler{r}
}

// ITEMS
func (h *BulkHandler) UploadItems(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")

	if err != nil {
		utils.ErrorResponse(w, 400, "file required")
		return
	}
	
	defer file.Close()

	rows, err := utils.ParseFile(file, header.Filename)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	userID := r.Context().Value("userID").(int64)

	result := h.repo.BulkInsertItems(rows, userID)

	utils.SuccessResponse(w, 200, result)
}

// STOCK
func (h *BulkHandler) UploadStock(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.ErrorResponse(w, 400, "file required")
		return
	}
	defer file.Close()

	rows, err := utils.ParseFile(file, header.Filename)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	result := h.repo.BulkInsertStock(rows)

	utils.SuccessResponse(w, 200, result)
}