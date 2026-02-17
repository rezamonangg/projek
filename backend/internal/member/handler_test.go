package member

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monachy/projek/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupMemberHandler(t *testing.T) (*Handler, *MockRepository) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	svc := NewService(repo)
	return NewHandler(svc), repo
}

func withMemberContext(r *http.Request, m *Member) *http.Request {
	ctx := context.WithValue(r.Context(), "member", m)
	return r.WithContext(ctx)
}

func withAdminMemberContext(r *http.Request) *http.Request {
	admin := testMember()
	admin.Role = RoleAdmin
	return withMemberContext(r, admin)
}

func withMemberChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestMemberHandler_List_Success(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	members := []Member{*testMember()}
	repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID, 20, 0).Return(members, nil)

	req := httptest.NewRequest("GET", "/members", nil)
	req = withMemberContext(req, testMember())

	rr := httptest.NewRecorder()
	handler.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestMemberHandler_List_Unauthorized(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	req := httptest.NewRequest("GET", "/members", nil)

	rr := httptest.NewRecorder()
	handler.List(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMemberHandler_Get_Success(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testMemberID).Return(testMember(), nil)

	req := httptest.NewRequest("GET", "/members/"+testMemberID.String(), nil)
	req = withMemberChiURLParams(req, map[string]string{"id": testMemberID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMemberHandler_Get_NotFound(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testMemberID).Return(nil, common.ErrNotFound)

	req := httptest.NewRequest("GET", "/members/"+testMemberID.String(), nil)
	req = withMemberChiURLParams(req, map[string]string{"id": testMemberID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMemberHandler_Get_InvalidID(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	req := httptest.NewRequest("GET", "/members/invalid", nil)
	req = withMemberChiURLParams(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMemberHandler_Invite_Success(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	body := `{"email":"new@example.com","role":"member"}`

	repo.EXPECT().GetByEmail(gomock.Any(), "new@example.com").Return(nil, common.ErrNotFound)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
	req = withAdminMemberContext(req)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Invite(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestMemberHandler_Invite_Unauthorized(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	body := `{"email":"new@example.com","role":"member"}`

	req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Invite(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMemberHandler_Invite_ForbiddenForNonAdmin(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	body := `{"email":"new@example.com","role":"member"}`

	req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
	req = withMemberContext(req, testMember())
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Invite(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestMemberHandler_Invite_DuplicateEmail(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	body := `{"email":"existing@example.com","role":"member"}`

	repo.EXPECT().GetByEmail(gomock.Any(), "existing@example.com").Return(testMember(), nil)

	req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(body))
	req = withAdminMemberContext(req)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Invite(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestMemberHandler_Invite_ValidationError(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	tests := []struct {
		name string
		body string
	}{
		{"missing email", `{"role":"member"}`},
		{"invalid email", `{"email":"invalid","role":"member"}`},
		{"missing role", `{"email":"test@example.com"}`},
		{"invalid role", `{"email":"test@example.com","role":"invalid"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/members/invite", strings.NewReader(tt.body))
			req = withAdminMemberContext(req)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Invite(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestMemberHandler_Update_Success(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	existing := testMember()
	repo.EXPECT().GetByID(gomock.Any(), testMemberID).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"first_name":"Updated","last_name":"Name"}`
	req := httptest.NewRequest("PUT", "/members/"+testMemberID.String(), strings.NewReader(body))
	req = withMemberChiURLParams(req, map[string]string{"id": testMemberID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMemberHandler_Update_NotFound(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testMemberID).Return(nil, common.ErrNotFound)

	body := `{"first_name":"Updated"}`
	req := httptest.NewRequest("PUT", "/members/"+testMemberID.String(), strings.NewReader(body))
	req = withMemberChiURLParams(req, map[string]string{"id": testMemberID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMemberHandler_Update_InvalidID(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	body := `{"first_name":"Updated"}`
	req := httptest.NewRequest("PUT", "/members/invalid", strings.NewReader(body))
	req = withMemberChiURLParams(req, map[string]string{"id": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMemberHandler_Delete_Success(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testMemberID).Return(nil)

	req := httptest.NewRequest("DELETE", "/members/"+testMemberID.String(), nil)
	req = withMemberChiURLParams(req, map[string]string{"id": testMemberID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestMemberHandler_Delete_NotFound(t *testing.T) {
	handler, repo := setupMemberHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testMemberID).Return(common.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/members/"+testMemberID.String(), nil)
	req = withMemberChiURLParams(req, map[string]string{"id": testMemberID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMemberHandler_Routes(t *testing.T) {
	handler, _ := setupMemberHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"POST", "/invite"},
		{"GET", "/{id}"},
		{"PUT", "/{id}"},
		{"DELETE", "/{id}"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var body *strings.Reader
			if tt.method == "POST" || tt.method == "PUT" {
				body = strings.NewReader(`{}`)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.method == "POST" || tt.method == "PUT" {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.NotEqual(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}
