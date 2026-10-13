package handlers

import (
	"encoding/json"
	"net/http"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type SalesHandler struct {
	salesRepo *repos.SalesRepo
}

func NewSalesHandler(r *repos.SalesRepo) *SalesHandler {
	return &SalesHandler{r}
}

type SaleRequest struct {
	Bill     models.Bill        `json:"bill"`
	Items    []models.BillItem  `json:"items"`
	Payments []models.Payment   `json:"payments"`
}

func (h *SalesHandler) CreateSale(w http.ResponseWriter, r *http.Request) {

	var req SaleRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, 400, "invalid request")
		return
	}

	req.Bill.CreatedBy = r.Context().Value("userID").(int64)

	err = h.salesRepo.CreateSale(req.Bill, req.Items, req.Payments)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 201, "Sale completed")
}