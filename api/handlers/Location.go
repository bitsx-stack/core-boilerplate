package handlers

import (
	"net/http"
	"strconv"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type LocationHandler struct {
	locationRepo *repos.LocationRepo
}

func NewLocationHandler(locationRepo *repos.LocationRepo) *LocationHandler {
	return &LocationHandler{locationRepo: locationRepo}
}

func (h *LocationHandler) GetAllLocations(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	model := q.Get("model")
	search := q.Get("search")
	country := q.Get("country")

	if country == "" {
		country = "TZ"
	}

	regionID := 0
	districtID := 0

	if v, err := strconv.Atoi(q.Get("region_id")); err == nil {
		regionID = v
	}
	if v, err := strconv.Atoi(q.Get("district_id")); err == nil {
		districtID = v
	}

	if model == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "model is required (continent, country, region, district, village)")
		return
	}

	data, err := h.locationRepo.GetAllLocations(model, search, country, regionID, districtID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"data": data,
	})
}