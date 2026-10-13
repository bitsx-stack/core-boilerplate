package handlers

import (
	"net/http"
	"strconv"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type ReportHandler struct {
	reportRepo *repos.ReportRepo
}

func NewReportHandler(rr *repos.ReportRepo) *ReportHandler {
	return &ReportHandler{reportRepo: rr}
}

// GET /api/reports/stock?item_id=&store_id=&start=&end=
func (h *ReportHandler) StockReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var itemID, storeID int64
	if v, err := strconv.ParseInt(q.Get("item_id"), 10, 64); err == nil {
		itemID = v
	}
	if v, err := strconv.ParseInt(q.Get("store_id"), 10, 64); err == nil {
		storeID = v
	}

	data, err := h.reportRepo.GetStockReport(itemID, storeID, q.Get("start"), q.Get("end"))
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	// Compute totals
	var totals repos.StockReportRow
	for _, row := range data {
		totals.OpenStock += row.OpenStock
		totals.PurchaseStock += row.PurchaseStock
		totals.TransferStock += row.TransferStock
		totals.AdjustStock += row.AdjustStock
		totals.SaleStock += row.SaleStock
		totals.CloseStock += row.CloseStock
		totals.PurchaseAmt += row.PurchaseAmt
		totals.SaleAmt += row.SaleAmt
		totals.BankAmt += row.BankAmt
		totals.CashAmt += row.CashAmt
		totals.MobileAmt += row.MobileAmt
		totals.Balance += row.Balance
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{
		"data":   data,
		"totals": totals,
	})
}