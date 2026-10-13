package handlers

import (
	"net/http"
	"strconv"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type DashboardHandler struct {
	dashboardRepo *repos.DashboardRepo
}

func NewDashboardHandler(dr *repos.DashboardRepo) *DashboardHandler {
	return &DashboardHandler{dashboardRepo: dr}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	var storeID int64
	if v, err := strconv.ParseInt(r.URL.Query().Get("store_id"), 10, 64); err == nil {
		storeID = v
	}

	data, err := h.dashboardRepo.GetDashboard(storeID)
	if err != nil {
		utils.ErrorResponse(w, 500, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, data)
}