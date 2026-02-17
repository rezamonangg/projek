package wiki

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

func setupWikiHandler(t *testing.T) (*Handler, *MockRepository) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	svc := NewService(repo)
	return NewHandler(svc), repo
}

func withWikiChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestWikiHandler_ListByProject_Success(t *testing.T) {
	handler, repo := setupWikiHandler(t)

	pages := []WikiPage{*testWikiPage()}
	repo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(pages, nil)

	req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/wiki", nil)
	req = withWikiChiURLParams(req, map[string]string{"projectId": testProjectID.String()})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestWikiHandler_Create_Success(t *testing.T) {
	handler, repo := setupWikiHandler(t)

	body := `{"title":"New Page","content":"Content"}`
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/wiki", strings.NewReader(body))
	req = withWikiChiURLParams(req, map[string]string{"projectId": testProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestWikiHandler_Create_MissingTitle(t *testing.T) {
	handler, _ := setupWikiHandler(t)

	body := `{"content":"Content"}`

	req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/wiki", strings.NewReader(body))
	req = withWikiChiURLParams(req, map[string]string{"projectId": testProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWikiHandler_Get_Success(t *testing.T) {
	handler, repo := setupWikiHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testPageID).Return(testWikiPage(), nil)

	req := httptest.NewRequest("GET", "/"+testPageID.String(), nil)
	req = withWikiChiURLParams(req, map[string]string{"id": testPageID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWikiHandler_Get_NotFound(t *testing.T) {
	handler, repo := setupWikiHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testPageID).Return(nil, common.ErrNotFound)

	req := httptest.NewRequest("GET", "/"+testPageID.String(), nil)
	req = withWikiChiURLParams(req, map[string]string{"id": testPageID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestWikiHandler_Update_Success(t *testing.T) {
	handler, repo := setupWikiHandler(t)

	existing := testWikiPage()
	repo.EXPECT().GetByID(gomock.Any(), testPageID).Return(existing, nil)
	repo.EXPECT().GetVersions(gomock.Any(), testPageID).Return([]WikiPageVersion{}, nil)
	repo.EXPECT().CreateVersion(gomock.Any(), gomock.Any()).Return(nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"title":"Updated Title","content":"Updated content"}`
	req := httptest.NewRequest("PUT", "/"+testPageID.String(), strings.NewReader(body))
	req = withWikiChiURLParams(req, map[string]string{"id": testPageID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWikiHandler_Delete_Success(t *testing.T) {
	handler, repo := setupWikiHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testPageID).Return(nil)

	req := httptest.NewRequest("DELETE", "/"+testPageID.String(), nil)
	req = withWikiChiURLParams(req, map[string]string{"id": testPageID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestWikiHandler_ProjectRoutes(t *testing.T) {
	handler, _ := setupWikiHandler(t)

	router := handler.ProjectRoutes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"POST", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var body *strings.Reader
			if tt.method == "POST" {
				body = strings.NewReader(`{}`)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.NotEqual(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}

func TestWikiHandler_Routes(t *testing.T) {
	handler, _ := setupWikiHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/{id}"},
		{"PUT", "/{id}"},
		{"DELETE", "/{id}"},
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
