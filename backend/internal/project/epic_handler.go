package project

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type EpicHandler struct {
	service *Service
}

func NewEpicHandler(service *Service) *EpicHandler {
	return &EpicHandler{service: service}
}

func (h *EpicHandler) ProjectRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	return r
}

func (h *EpicHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
	return r
}

func (h *EpicHandler) List(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	epics, err := h.service.GetEpics(r.Context(), projectID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, epics)
}

func (h *EpicHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	var req CreateEpicInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	epic, err := h.service.CreateEpic(r.Context(), projectID, req)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, epic)
}

func (h *EpicHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid epic id")
		return
	}

	var req CreateEpicInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	epic, err := h.service.UpdateEpic(r.Context(), id, req)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "epic not found")
		return
	}

	common.Success(w, http.StatusOK, epic)
}

func (h *EpicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid epic id")
		return
	}

	if err := h.service.DeleteEpic(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "epic not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}
