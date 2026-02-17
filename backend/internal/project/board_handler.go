package project

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type BoardHandler struct {
	service *Service
}

func NewBoardHandler(service *Service) *BoardHandler {
	return &BoardHandler{service: service}
}

func (h *BoardHandler) ProjectRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByProject)
	r.Post("/", h.Create)
	return r
}

func (h *BoardHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Delete("/", h.Delete)
	})
	return r
}

func (h *BoardHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	boards, err := h.service.GetBoards(r.Context(), projectID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, boards)
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	var req CreateBoardInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	board, err := h.service.CreateBoard(r.Context(), projectID, req)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, board)
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid board id")
		return
	}

	board, err := h.service.GetBoardByID(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "board not found")
		return
	}

	common.Success(w, http.StatusOK, board)
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid board id")
		return
	}

	if err := h.service.DeleteBoard(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "board not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}
