package project

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/member"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testMemberID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

func testMember() *member.Member {
	return &member.Member{
		ID:          testMemberID,
		CommunityID: testCommunityID,
		Email:       "test@example.com",
		Role:        "member",
		IsActive:    true,
	}
}

func setupHandler(t *testing.T) (*Handler, *MockRepository, *MockEpicRepository, *MockBoardRepository) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	epicRepo := NewMockEpicRepository(ctrl)
	boardRepo := NewMockBoardRepository(ctrl)
	svc := NewService(repo, epicRepo, boardRepo, nil, nil)
	return NewHandler(svc), repo, epicRepo, boardRepo
}

func withMember(r *http.Request, m *member.Member) *http.Request {
	ctx := context.WithValue(r.Context(), "member", m)
	return r.WithContext(ctx)
}

func withChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestHandler_Create_Success(t *testing.T) {
	handler, repo, _, boardRepo := setupHandler(t)

	body := `{"name":"New Project","description":"Desc","key":"NEW"}`
	req := httptest.NewRequest("POST", "/projects", strings.NewReader(body))
	req = withMember(req, testMember())
	req.Header.Set("Content-Type", "application/json")

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestHandler_Create_ValidationError(t *testing.T) {
	handler, _, _, _ := setupHandler(t)

	tests := []struct {
		name string
		body string
	}{
		{"empty name", `{"name":"","description":"Desc","key":"NEW"}`},
		{"missing key", `{"name":"Project","description":"Desc"}`},
		{"key too short", `{"name":"Project","key":"A"}`},
		{"key too long", `{"name":"Project","key":"TOOLONGKEY1"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/projects", strings.NewReader(tt.body))
			req = withMember(req, testMember())
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Create(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestHandler_Create_Unauthorized(t *testing.T) {
	handler, _, _, _ := setupHandler(t)

	body := `{"name":"New Project","description":"Desc","key":"NEW"}`
	req := httptest.NewRequest("POST", "/projects", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandler_List_Success(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	projects := []Project{*testProject()}
	repo.EXPECT().GetByCommunity(gomock.Any(), testCommunityID, 20, 0).Return(projects, nil)

	req := httptest.NewRequest("GET", "/projects", nil)
	req = withMember(req, testMember())

	rr := httptest.NewRecorder()
	handler.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestHandler_List_Unauthorized(t *testing.T) {
	handler, _, _, _ := setupHandler(t)

	req := httptest.NewRequest("GET", "/projects", nil)

	rr := httptest.NewRecorder()
	handler.List(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(testProject(), nil)

	req := httptest.NewRequest("GET", "/projects/"+testProjectID.String(), nil)
	req = withChiURLParam(req, "id", testProjectID.String())

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(nil, common.ErrNotFound)

	req := httptest.NewRequest("GET", "/projects/"+testProjectID.String(), nil)
	req = withChiURLParam(req, "id", testProjectID.String())

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	handler, _, _, _ := setupHandler(t)

	req := httptest.NewRequest("GET", "/projects/invalid", nil)
	req = withChiURLParam(req, "id", "invalid")

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	existing := testProject()
	repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest("PUT", "/projects/"+testProjectID.String(), strings.NewReader(body))
	req = withChiURLParam(req, "id", testProjectID.String())
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_Update_NotFound(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testProjectID).Return(nil, common.ErrNotFound)

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest("PUT", "/projects/"+testProjectID.String(), strings.NewReader(body))
	req = withChiURLParam(req, "id", testProjectID.String())
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testProjectID).Return(nil)

	req := httptest.NewRequest("DELETE", "/projects/"+testProjectID.String(), nil)
	req = withChiURLParam(req, "id", testProjectID.String())

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	handler, repo, _, _ := setupHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testProjectID).Return(common.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/projects/"+testProjectID.String(), nil)
	req = withChiURLParam(req, "id", testProjectID.String())

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandler_Routes(t *testing.T) {
	handler, _, _, _ := setupHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"POST", "/"},
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

func TestQueryInt(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		key          string
		defaultValue int
		expected     int
	}{
		{"no query param", "", "page", 1, 1},
		{"valid int", "page=5", "page", 1, 5},
		{"invalid int", "page=abc", "page", 1, 1},
		{"negative int", "page=-1", "page", 1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/?"+tt.query, nil)
			result := queryInt(req, tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaginatedResponse(t *testing.T) {
	projects := []Project{*testProject(), *testProject()}
	response := PaginatedResponse[Project]{
		Items:      projects,
		Total:      2,
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(response)
	require.NoError(t, err)

	var decoded PaginatedResponse[Project]
	err = json.NewDecoder(buf).Decode(&decoded)
	require.NoError(t, err)

	assert.Len(t, decoded.Items, 2)
	assert.Equal(t, 2, decoded.Total)
	assert.Equal(t, 1, decoded.Page)
	assert.Equal(t, 20, decoded.PageSize)
}
