package handlers

import (
	"net/http"
	"strconv"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type StockHandler struct {
	stockRepo *repos.StockRepo
	transRepo *repos.TransactionRepo
}

func NewStockHandler(s *repos.StockRepo, t *repos.TransactionRepo) *StockHandler {
	return &StockHandler{s, t}
}

func (h *StockHandler) AddStock(w http.ResponseWriter, r *http.Request) {

	itemID, _ := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	storeID, _ := strconv.ParseInt(r.FormValue("store_id"), 10, 64)
	qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	

	err := h.stockRepo.AddStock(itemID, storeID, qty)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}


	utils.SuccessResponse(w, 200, "Stock added")
}

func (h *StockHandler) ReduceStock(w http.ResponseWriter, r *http.Request) {

	itemID, _ := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	storeID, _ := strconv.ParseInt(r.FormValue("store_id"), 10, 64)
	qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	

	err := h.stockRepo.ReduceStock(itemID, storeID, qty)
	if err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}


	utils.SuccessResponse(w, 200, "Stock reduced")
}