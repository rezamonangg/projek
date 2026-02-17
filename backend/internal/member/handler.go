package member

import (
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/invite", h.Invite)
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

type InviteRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  Role   `json:"role" validate:"required,oneof=admin member"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func getCurrentMember(r *http.Request) *Member {
	if v := r.Context().Value("member"); v != nil {
		if m, ok := v.(*Member); ok {
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

	members, err := h.service.GetByCommunity(r.Context(), m.CommunityID, pageSize, offset)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response := PaginatedResponse[Member]{
		Items:    members,
		Total:    len(members),
		Page:     page,
		PageSize: pageSize,
	}
	if len(members) > 0 {
		response.TotalPages = (len(members) + pageSize - 1) / pageSize
	}

	common.Success(w, http.StatusOK, response)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid member id")
		return
	}

	member, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "member not found")
		return
	}

	common.Success(w, http.StatusOK, member)
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	m := getCurrentMember(r)
	if m == nil {
		common.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	if m.Role != RoleAdmin {
		common.Error(w, http.StatusForbidden, "FORBIDDEN", "only admins can invite members")
		return
	}

	var req InviteRequest
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	existing, _ := h.service.GetByEmail(r.Context(), req.Email)
	if existing != nil {
		common.Error(w, http.StatusConflict, "CONFLICT", "email already exists")
		return
	}

	input := CreateMemberInput{
		Email:    req.Email,
		Password: "temp_password_123",
		Role:     req.Role,
	}

	newMember, err := h.service.Create(r.Context(), m.CommunityID, input)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, map[string]string{"id": newMember.ID.String()})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid member id")
		return
	}

	var req UpdateProfileRequest
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	input := UpdateMemberInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	member, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.Error(w, http.StatusNotFound, "NOT_FOUND", "member not found")
			return
		}
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, member)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid member id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "member not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}
