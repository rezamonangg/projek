package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Get("/stats", h.GetStats)
	r.Get("/settings", h.GetSettings)
	r.Put("/settings", h.UpdateSettings)
	return r
}

type UpdateSettingsRequest struct {
	AllowMemberRegistration  bool `json:"allow_member_registration"`
	RequireEmailVerification bool `json:"require_email_verification"`
}

func getCurrentMember(r *http.Request) *member.Member {
	if v := r.Context().Value("member"); v != nil {
		if m, ok := v.(*member.Member); ok {
			return m
		}
	}
	return nil
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	m := getCurrentMember(r)
	if m == nil {
		common.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	stats, err := h.service.GetDashboardStats(r.Context(), m.CommunityID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, stats)
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	m := getCurrentMember(r)
	if m == nil {
		common.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	settings, err := h.service.GetSettings(r.Context(), m.CommunityID)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "settings not found")
		return
	}

	common.Success(w, http.StatusOK, settings)
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	m := getCurrentMember(r)
	if m == nil {
		common.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	if m.Role != member.RoleAdmin {
		common.Error(w, http.StatusForbidden, "FORBIDDEN", "only admins can update settings")
		return
	}

	var req UpdateSettingsRequest
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	settings, err := h.service.UpdateSettings(r.Context(), m.CommunityID, req.AllowMemberRegistration, req.RequireEmailVerification)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, settings)
}
