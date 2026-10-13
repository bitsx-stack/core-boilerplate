package handlers

import (
	"net/http"
	"strings"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type RegionHandler struct {
	regionRepo *repos.RegionRepo
}

func NewRegionHandler(regionRepo *repos.RegionRepo) *RegionHandler {

	return &RegionHandler{
		regionRepo: regionRepo,
	}
}

func (h *RegionHandler) GetAllRegions(response http.ResponseWriter, request *http.Request) {

	countryStr := strings.TrimSpace(request.URL.Query().Get("country"))
	regions, err := h.regionRepo.GetAllRegions(countryStr)
	if err != nil {
		utils.ErrorResponse(response, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(response, http.StatusOK, regions)

}
