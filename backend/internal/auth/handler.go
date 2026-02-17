package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/member"
)

type Handler struct {
	authService   *AuthService
	memberService *member.Service
}

func NewHandler(authService *AuthService, memberService *member.Service) *Handler {
	return &Handler{
		authService:   authService,
		memberService: memberService,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	r.Get("/me", h.Me)
	return r
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	SessionID string         `json:"session_id"`
	Member    *member.Member `json:"member"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid email or password")
		return
	}

	session, m, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		common.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	common.Success(w, http.StatusOK, LoginResponse{
		SessionID: session.ID,
		Member:    m,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID != "" {
		h.authService.Logout(r.Context(), sessionID)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	common.Success(w, http.StatusOK, nil)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	member := getCurrentMember(r)
	if member == nil {
		common.Error(w, http.StatusUnauthorized, "NOT_AUTHENTICATED", "not authenticated")
		return
	}
	common.Success(w, http.StatusOK, member)
}

func getSessionID(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getCurrentMember(r *http.Request) *member.Member {
	if v := r.Context().Value("member"); v != nil {
		if m, ok := v.(*member.Member); ok {
			return m
		}
	}
	return nil
}
