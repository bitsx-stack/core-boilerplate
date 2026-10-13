package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

// ══════════════════════════════════════════════
// BILL HANDLER
// ══════════════════════════════════════════════

type BillHandler struct {
	billRepo    *repos.BillRepo
	paymentRepo *repos.PaymentRepo
	stockRepo   *repos.StockRepo
}

func NewBillHandler(br *repos.BillRepo, pr *repos.PaymentRepo, sr *repos.StockRepo) *BillHandler {
	return &BillHandler{billRepo: br, paymentRepo: pr, stockRepo: sr}
}

// POST /api/bills/create
func (h *BillHandler) CreateBill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PeopleID int64  `json:"people_id"`
		StoreID  int64  `json:"store_id"`
		Type     string `json:"type"`
		Items    []struct {
			ItemID         int64   `json:"item_id"`
			ItemPrice      float64 `json:"item_price"`
			ItemQuantity   float64 `json:"item_quantity"`
			DiscountAmount float64 `json:"discount_amount"`
			DiscountRate   float64 `json:"discount_rate"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	if req.PeopleID == 0 {
		utils.ErrorResponse(w, 400, "people_id is required")
		return
	}
	if req.StoreID == 0 {
		utils.ErrorResponse(w, 400, "store_id is required")
		return
	}
	if req.Type == "" {
		utils.ErrorResponse(w, 400, "type is required")
		return
	}
	if len(req.Items) == 0 {
		utils.ErrorResponse(w, 400, "items required")
		return
	}

	userID := r.Context().Value("userID").(int64)
	code := fmt.Sprintf("%s-%d", strings.ToUpper(req.Type[:3]), time.Now().UnixMilli()%100000)

	bill := models.Bill{
		Code: code, PeopleID: req.PeopleID, StoreID: req.StoreID,
		Type: req.Type, CreatedBy: userID,
	}

	var items []models.BillItem
	for _, item := range req.Items {
		items = append(items, models.BillItem{
			ItemID: item.ItemID, ItemPrice: item.ItemPrice,
			ItemQuantity: item.ItemQuantity, DiscountAmount: item.DiscountAmount,
			DiscountRate: item.DiscountRate,
		})
	}

	billID, err := h.billRepo.CreateBill(bill, items, h.stockRepo)
	if err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	utils.SuccessResponse(w, 201, map[string]interface{}{
		"message": "Bill created", "bill_id": billID, "code": code,
	})
}



// GET /api/bills
func (h *BillHandler) GetAllBills(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	billType := q.Get("type")
	status := q.Get("status")
	store := q.Get("store")

	var peopleID int64
	if v, err := strconv.ParseInt(q.Get("people_id"), 10, 64); err == nil {
		peopleID = v
	}

	limit := 10
	offset := 0
	if v, err := strconv.Atoi(q.Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(q.Get("offset")); err == nil && v >= 0 {
		offset = v
	}

	data, total, err := h.billRepo.GetAllBills(search, billType, status, store, peopleID, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{
		"data":  data,
		"total": total,
	})
}

// GET /api/bills/single/{id}
func (h *BillHandler) GetSingleBill(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/bills/single/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid bill ID")
		return
	}

	bill, err := h.billRepo.GetSingleBill(id)
	if err != nil {
		utils.ErrorResponse(w, 404, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, bill)
}

// PUT /api/bills/update/{id}
func (h *BillHandler) UpdateBillItems(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/bills/update/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid bill ID")
		return
	}

	var req struct {
		Items []struct {
			ItemID         int64   `json:"item_id"`
			ItemPrice      float64 `json:"item_price"`
			ItemQuantity   float64 `json:"item_quantity"`
			DiscountAmount float64 `json:"discount_amount"`
			DiscountRate   float64 `json:"discount_rate"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	if len(req.Items) == 0 {
		utils.ErrorResponse(w, 400, "at least one item is required")
		return
	}

	var items []models.BillItem
	for _, item := range req.Items {
		items = append(items, models.BillItem{
			ItemID:         item.ItemID,
			ItemPrice:      item.ItemPrice,
			ItemQuantity:   item.ItemQuantity,
			DiscountAmount: item.DiscountAmount,
			DiscountRate:   item.DiscountRate,
		})
	}

	err = h.billRepo.UpdateBillItems(id, items)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "Bill updated")
}

// POST /api/bills/pay/{id}
func (h *BillHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/bills/pay/")
	billID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid bill ID")
		return
	}

	var req struct {
		Type          string  `json:"type"`
		Category      string  `json:"category"`
		PaymentAmount float64 `json:"payment_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	if req.PaymentAmount <= 0 {
		utils.ErrorResponse(w, 400, "payment_amount must be greater than 0")
		return
	}
	if req.Type == "" {
		utils.ErrorResponse(w, 400, "type is required (CASH, MOBILE, BANK)")
		return
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		utils.ErrorResponse(w, 401, "unauthorized")
		return
	}

	// Get bill to find people_id
	bill, err := h.billRepo.GetSingleBill(billID)
	if err != nil {
		utils.ErrorResponse(w, 404, err.Error())
		return
	}

	payment := models.Payment{
		BillID:        billID,
		PeopleID:      bill.PeopleID,
		Type:          req.Type,
		Category:      req.Category,
		PaymentAmount: req.PaymentAmount,
		CreatedBy:     userID,
	}

	err = h.paymentRepo.CreatePayment(payment, h.billRepo)
	if err != nil {
		utils.ErrorResponse(w, 400, err.Error())
		return
	}

	utils.SuccessResponse(w, 201, "Payment recorded")
}

// GET /api/bills/payments/{id}
func (h *BillHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/bills/payments/")
	billID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, 400, "Invalid bill ID")
		return
	}

	payments, err := h.paymentRepo.GetPaymentsByBill(billID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, map[string]interface{}{"data": payments})
}