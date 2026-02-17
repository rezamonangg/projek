package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/member"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testCommunityID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var testAdminMemberID = uuid.MustParse("22222222-2222-2222-2222-222222222222")

func testAdminMember() *member.Member {
	return &member.Member{
		ID:          testAdminMemberID,
		CommunityID: testCommunityID,
		Email:       "admin@example.com",
		Role:        member.RoleAdmin,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func testRegularMember() *member.Member {
	return &member.Member{
		ID:          uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		CommunityID: testCommunityID,
		Email:       "member@example.com",
		Role:        member.RoleMember,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func testSettings() *CommunitySettings {
	return &CommunitySettings{
		ID:                       uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		CommunityID:              testCommunityID,
		AllowMemberRegistration:  true,
		RequireEmailVerification: false,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}
}

func setupAdminHandler(t *testing.T) (*Handler, *MockSettingsRepository, *MockStatsRepository) {
	ctrl := gomock.NewController(t)
	settingsRepo := NewMockSettingsRepository(ctrl)
	statsRepo := NewMockStatsRepository(ctrl)
	svc := NewService(settingsRepo, statsRepo)
	return NewHandler(svc), settingsRepo, statsRepo
}

func withMemberContext(r *http.Request, m *member.Member) *http.Request {
	ctx := context.WithValue(r.Context(), "member", m)
	return r.WithContext(ctx)
}

func withAdminChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestAdminHandler_GetStats_Success(t *testing.T) {
	handler, _, statsRepo := setupAdminHandler(t)

	statsRepo.EXPECT().GetDashboardStats(gomock.Any(), testCommunityID).Return(&DashboardStats{}, nil)

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	req = withMemberContext(req, testAdminMember())

	rr := httptest.NewRecorder()
	handler.GetStats(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAdminHandler_GetStats_Unauthorized(t *testing.T) {
	handler, _, _ := setupAdminHandler(t)

	req := httptest.NewRequest("GET", "/admin/stats", nil)

	rr := httptest.NewRecorder()
	handler.GetStats(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAdminHandler_GetSettings_Success(t *testing.T) {
	handler, repo, _ := setupAdminHandler(t)

	repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID).Return(testSettings(), nil)

	req := httptest.NewRequest("GET", "/admin/settings", nil)
	req = withMemberContext(req, testAdminMember())

	rr := httptest.NewRecorder()
	handler.GetSettings(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestAdminHandler_GetSettings_Unauthorized(t *testing.T) {
	handler, _, _ := setupAdminHandler(t)

	req := httptest.NewRequest("GET", "/admin/settings", nil)

	rr := httptest.NewRecorder()
	handler.GetSettings(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAdminHandler_UpdateSettings_Success(t *testing.T) {
	handler, repo, _ := setupAdminHandler(t)

	repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID).Return(testSettings(), nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"allow_member_registration":false,"require_email_verification":true}`
	req := httptest.NewRequest("PUT", "/admin/settings", strings.NewReader(body))
	req = withMemberContext(req, testAdminMember())
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.UpdateSettings(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAdminHandler_UpdateSettings_ForbiddenForNonAdmin(t *testing.T) {
	handler, _, _ := setupAdminHandler(t)

	body := `{"allow_member_registration":false}`
	req := httptest.NewRequest("PUT", "/admin/settings", strings.NewReader(body))
	req = withMemberContext(req, testRegularMember())
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.UpdateSettings(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAdminHandler_UpdateSettings_Unauthorized(t *testing.T) {
	handler, _, _ := setupAdminHandler(t)

	body := `{"allow_member_registration":false}`
	req := httptest.NewRequest("PUT", "/admin/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.UpdateSettings(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAdminHandler_Routes(t *testing.T) {
	handler, _, _ := setupAdminHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/stats"},
		{"GET", "/settings"},
		{"PUT", "/settings"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var body *strings.Reader
			if tt.method == "PUT" {
				body = strings.NewReader(`{}`)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.method == "PUT" {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.NotEqual(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}
