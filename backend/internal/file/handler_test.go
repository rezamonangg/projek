package file

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monachy/projek/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupFileHandler(t *testing.T) (*Handler, *MockRepository, Storage) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	storage := NewLocalStorage("/tmp", "/uploads")
	svc := NewService(repo, storage)
	return NewHandler(svc), repo, storage
}

func withFileChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestFileHandler_ListByProject_Success(t *testing.T) {
	handler, repo, _ := setupFileHandler(t)

	attachments := []FileAttachment{*testFileAttachment()}
	repo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(attachments, nil)

	req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/files", nil)
	req = withFileChiURLParams(req, map[string]string{"projectId": testProjectID.String()})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestFileHandler_ListByProject_InvalidProjectID(t *testing.T) {
	handler, _, _ := setupFileHandler(t)

	req := httptest.NewRequest("GET", "/invalid/files", nil)
	req = withFileChiURLParams(req, map[string]string{"projectId": "invalid"})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFileHandler_Get_Success(t *testing.T) {
	handler, repo, _ := setupFileHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testFileID).Return(testFileAttachment(), nil)

	req := httptest.NewRequest("GET", "/"+testFileID.String(), nil)
	req = withFileChiURLParams(req, map[string]string{"id": testFileID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFileHandler_Get_NotFound(t *testing.T) {
	handler, repo, _ := setupFileHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testFileID).Return(nil, common.ErrNotFound)

	req := httptest.NewRequest("GET", "/"+testFileID.String(), nil)
	req = withFileChiURLParams(req, map[string]string{"id": testFileID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestFileHandler_Delete_Success(t *testing.T) {
	handler, repo, _ := setupFileHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testFileID).Return(testFileAttachment(), nil)
	repo.EXPECT().Delete(gomock.Any(), testFileID).Return(nil)

	req := httptest.NewRequest("DELETE", "/"+testFileID.String(), nil)
	req = withFileChiURLParams(req, map[string]string{"id": testFileID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestFileHandler_Routes(t *testing.T) {
	handler, _, _ := setupFileHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"POST", "/"},
		{"GET", "/{id}"},
		{"DELETE", "/{id}"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.NotEqual(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}
