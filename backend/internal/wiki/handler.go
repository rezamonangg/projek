package wiki

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ProjectRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByProject)
	r.Post("/", h.Create)
	return r
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
	r.Get("/{id}/versions", h.GetVersions)
	return r
}

func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	pages, err := h.service.GetByProject(r.Context(), projectID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, pages)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	var req CreateWikiPageInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	page, err := h.service.Create(r.Context(), projectID, req)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, page)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid wiki page id")
		return
	}

	page, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "wiki page not found")
		return
	}

	common.Success(w, http.StatusOK, page)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid wiki page id")
		return
	}

	var req UpdateWikiPageInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	var userID uuid.UUID
	if member := getCurrentMember(r); member != nil {
		userID = member.ID
	}

	page, err := h.service.Update(r.Context(), id, req, userID)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "wiki page not found")
		return
	}

	common.Success(w, http.StatusOK, page)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid wiki page id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "wiki page not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}

func (h *Handler) GetVersions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid wiki page id")
		return
	}

	versions, err := h.service.GetVersions(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, versions)
}

type Member struct {
	ID uuid.UUID `json:"id"`
}

func getCurrentMember(r *http.Request) *Member {
	if v := r.Context().Value("member"); v != nil {
		if m, ok := v.(*Member); ok {
			return m
		}
	}
	return nil
}
