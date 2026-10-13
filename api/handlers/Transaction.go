package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type TransactionHandler struct {
	transactionRepo *repos.TransactionRepo
	billRepo        *repos.BillRepo
	stockRepo       *repos.StockRepo
}

func NewTransactionHandler(
	tr *repos.TransactionRepo,
	br *repos.BillRepo,
	sr *repos.StockRepo,
) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: tr,
		billRepo:        br,
		stockRepo:       sr,
	}
}


// func (h *TransactionHandler) CompleteTransaction(w http.ResponseWriter, r *http.Request) {

// 	txIDStr := r.PathValue("id")
// 	if txIDStr == "" {
// 		utils.ErrorResponse(w, 400, "transaction id is required")
// 		return
// 	}

// 	txID, err := strconv.ParseInt(txIDStr, 10, 64)
// 	if err != nil {
// 		utils.ErrorResponse(w, 400, "invalid transaction id")
// 		return
// 	}

// 	err = h.transactionRepo.CompleteTransaction(
// 		txID,
// 		h.stockRepo,
// 	)

// 	if err != nil {
// 		utils.ErrorResponse(w, 500, err.Error())
// 		return
// 	}

// 	utils.SuccessResponse(w, 200, "Transaction completed successfully")
// }

func (h *TransactionHandler) CompleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.transactionRepo.CompleteTransaction(id, h.stockRepo)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Transaction completed")
}


func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StoreID  int64  `json:"store_id"`
		Category string `json:"category"`
		Items    []struct {
			ItemID   int64   `json:"item_id"`
			Quantity float64 `json:"quantity"`
			Type     string  `json:"type"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.StoreID == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "store_id is required")
		return
	}
	if req.Category == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "category is required")
		return
	}
	if len(req.Items) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "No items provided")
		return
	}

	createdBy := r.Context().Value("userID").(int64)

	for _, item := range req.Items {
		if item.ItemID == 0 {
			utils.ErrorResponse(w, http.StatusBadRequest, "item_id is required")
			return
		}

		t := models.Transaction{
			ItemID:    item.ItemID,
			StoreID:   req.StoreID,
			Type:      item.Type,
			Category:  req.Category,
			Quantity:  item.Quantity,
			Status:    "pending",
			CreatedBy: createdBy,
		}

		if err := h.transactionRepo.CreateTransaction(t); err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.SuccessResponse(w, http.StatusCreated, "Transactions created successfully")
}


func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	search := query.Get("search")
	store := query.Get("store")
	typeFilter := query.Get("type")
	statusFilter := query.Get("status")
	categoryFilter := query.Get("category")

	limit := 10
	offset := 0

	if v, err := strconv.Atoi(query.Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(query.Get("offset")); err == nil && v >= 0 {
		offset = v
	}

	data, total, err := h.transactionRepo.GetAllTransactions(
		search, store, typeFilter, statusFilter, categoryFilter,
		limit, offset,
	)

	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{
		"data":  data,
		"total": total,
	})
}