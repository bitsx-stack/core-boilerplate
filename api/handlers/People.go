package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"bitpos.com/m/models"
	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type PeopleHandler struct {
	peopleRepo *repos.PeopleRepo
}

func NewPeopleHandler(peopleReople *repos.PeopleRepo) *PeopleHandler {
	return &PeopleHandler{
		peopleRepo: peopleReople,
	}
}

func (h *PeopleHandler) GetAllPeople(w http.ResponseWriter, r *http.Request) {

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	itemType := strings.TrimSpace(r.URL.Query().Get("item_type"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))

	limit := 10
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if oStr := r.URL.Query().Get("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o >= 0 {
			offset = o
		}
	}

	fmt.Printf("Query Params => search: %q, type: %q, category: %q, limit: %d, offset: %d\n",
		search, itemType, category, limit, offset)

	data, total, err := h.peopleRepo.GetAllPeople(search, itemType, category, limit, offset)
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

func (h *PeopleHandler) CreatePeople(w http.ResponseWriter, r *http.Request) {

	tags := r.FormValue("tags")

	req := models.People{
		Type:     r.FormValue("type"),
		Category: r.FormValue("category"),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
		Email:    r.FormValue("email"),
		Location: r.FormValue("location"),
		Tin:      r.FormValue("tin"),
		Vrn:      r.FormValue("vrn"),
		Nin:      r.FormValue("nin"),
		PoBox:    r.FormValue("po_box"),
		Tags:     &tags,
		Status:   r.FormValue("status") == "true",
	}

	if req.Type == "" || req.Category == "" || req.Location == "" || req.Name == "" || req.Phone == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name, Type, Category, Location and Phone are required")
		return
	}

	req.CreatedBy = r.Context().Value("userID").(int64)
	req.Name = strings.ToUpper(req.Name)

	logoFileName, err := utils.UploadFileHandler(w, r, "logo")
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "File upload failed: "+err.Error())
		return
	}

	req.LogoURL = logoFileName

	err = h.peopleRepo.CreatePeople(req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "People Created Successfully!")

}

func (h *PeopleHandler) GetSinglePeople(response http.ResponseWriter, request *http.Request) {

	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(response, http.StatusBadRequest, "Invalid Request")
		return
	}

	people, err := h.peopleRepo.GetSinglePeople(int64(id))
	if err != nil {
		utils.ErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	if people == nil {
		utils.ErrorResponse(response, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(response, http.StatusOK, people)

}


func (h *PeopleHandler) UpdatePeople(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/people/update/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid people ID")
		return
	}

	tags := r.FormValue("tags")

	req := models.People{
		Type:     r.FormValue("type"),
		Category: r.FormValue("category"),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
		Email:    r.FormValue("email"),
		Location: r.FormValue("location"),
		Tin:      r.FormValue("tin"),
		Vrn:      r.FormValue("vrn"),
		Nin:      r.FormValue("nin"),
		PoBox:    r.FormValue("po_box"),
		Tags:     &tags,
	}

	if req.Type == "" || req.Category == "" || req.Location == "" || req.Name == "" || req.Phone == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name, Type, Category, Location and Phone are required")
		return
	}

	req.Name = strings.ToUpper(req.Name)

	// Handle optional logo upload
	logoFileName, err := utils.UploadFileHandler(w, r, "logo")
	if err == nil && logoFileName != "" {
		req.LogoURL = logoFileName
	}

	err = h.peopleRepo.UpdatePeople(id, req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "People Updated Successfully!")
}

func (h *PeopleHandler) DeletePeople(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/people/delete/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid people ID")
		return
	}

	deletedBy := r.Context().Value("userID").(int64)

	err = h.peopleRepo.DeletePeople(id, deletedBy)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "People Deleted Successfully!")
}
