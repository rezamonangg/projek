package project

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/member"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
	return r
}

type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

func getCurrentMember(r *http.Request) *member.Member {
	if v := r.Context().Value("member"); v != nil {
		if m, ok := v.(*member.Member); ok {
			return m
		}
	}
	return nil
}

func queryInt(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	m := getCurrentMember(r)
	if m == nil {
		common.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	page := queryInt(r, "page", 1)
	pageSize := queryInt(r, "page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	projects, err := h.service.GetByCommunity(r.Context(), m.CommunityID, pageSize, offset)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response := PaginatedResponse[Project]{
		Items:    projects,
		Total:    len(projects),
		Page:     page,
		PageSize: pageSize,
	}
	if len(projects) > 0 {
		response.TotalPages = (len(projects) + pageSize - 1) / pageSize
	}

	common.Success(w, http.StatusOK, response)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	m := getCurrentMember(r)
	if m == nil {
		common.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	var req CreateProjectInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	project, err := h.service.Create(r.Context(), m.CommunityID, req)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, project)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid project id")
		return
	}

	project, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "project not found")
		return
	}

	common.Success(w, http.StatusOK, project)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid project id")
		return
	}

	var req UpdateProjectInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	project, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "project not found")
		return
	}

	common.Success(w, http.StatusOK, project)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid project id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "project not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}
