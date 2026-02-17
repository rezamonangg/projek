package project

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testEpicID = uuid.MustParse("44444444-4444-4444-4444-444444444444")

func testEpic() *Epic {
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now.Add(7 * 24 * time.Hour)
	return &Epic{
		ID:          testEpicID,
		ProjectID:   testProjectID,
		Name:        "Test Epic",
		Description: "Epic description",
		StartDate:   &startDate,
		EndDate:     &endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func setupEpicHandler(t *testing.T) (*EpicHandler, *MockEpicRepository) {
	ctrl := gomock.NewController(t)
	epicRepo := NewMockEpicRepository(ctrl)
	svc := NewService(nil, epicRepo, nil, nil, nil)
	return NewEpicHandler(svc), epicRepo
}

func withChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestEpicHandler_List_Success(t *testing.T) {
	handler, epicRepo := setupEpicHandler(t)

	epics := []Epic{*testEpic()}
	epicRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(epics, nil)

	req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/epics", nil)
	req = withChiURLParams(req, map[string]string{"projectId": testProjectID.String()})

	rr := httptest.NewRecorder()
	handler.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestEpicHandler_List_InvalidProjectID(t *testing.T) {
	handler, _ := setupEpicHandler(t)

	req := httptest.NewRequest("GET", "/invalid/epics", nil)
	req = withChiURLParams(req, map[string]string{"projectId": "invalid"})

	rr := httptest.NewRecorder()
	handler.List(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEpicHandler_Create_Success(t *testing.T) {
	handler, epicRepo := setupEpicHandler(t)

	body := `{"name":"New Epic","description":"Description"}`
	epicRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/epics", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"projectId": testProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestEpicHandler_Create_MissingName(t *testing.T) {
	handler, _ := setupEpicHandler(t)

	body := `{"description":"Description"}`

	req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/epics", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"projectId": testProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEpicHandler_Create_InvalidProjectID(t *testing.T) {
	handler, _ := setupEpicHandler(t)

	body := `{"name":"New Epic","description":"Description"}`

	req := httptest.NewRequest("POST", "/invalid/epics", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"projectId": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEpicHandler_Update_Success(t *testing.T) {
	handler, epicRepo := setupEpicHandler(t)

	existing := testEpic()
	epicRepo.EXPECT().GetByID(gomock.Any(), testEpicID).Return(existing, nil)
	epicRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"name":"Updated Epic","description":"Updated"}`
	req := httptest.NewRequest("PUT", "/"+testEpicID.String(), strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"id": testEpicID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEpicHandler_Update_NotFound(t *testing.T) {
	handler, epicRepo := setupEpicHandler(t)

	epicRepo.EXPECT().GetByID(gomock.Any(), testEpicID).Return(nil, common.ErrNotFound)

	body := `{"name":"Updated Epic","description":"Updated"}`
	req := httptest.NewRequest("PUT", "/"+testEpicID.String(), strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"id": testEpicID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestEpicHandler_Update_InvalidID(t *testing.T) {
	handler, _ := setupEpicHandler(t)

	body := `{"name":"Updated Epic","description":"Updated"}`
	req := httptest.NewRequest("PUT", "/invalid", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"id": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEpicHandler_Delete_Success(t *testing.T) {
	handler, epicRepo := setupEpicHandler(t)

	epicRepo.EXPECT().Delete(gomock.Any(), testEpicID).Return(nil)

	req := httptest.NewRequest("DELETE", "/"+testEpicID.String(), nil)
	req = withChiURLParams(req, map[string]string{"id": testEpicID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestEpicHandler_Delete_NotFound(t *testing.T) {
	handler, epicRepo := setupEpicHandler(t)

	epicRepo.EXPECT().Delete(gomock.Any(), testEpicID).Return(common.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/"+testEpicID.String(), nil)
	req = withChiURLParams(req, map[string]string{"id": testEpicID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestEpicHandler_Delete_InvalidID(t *testing.T) {
	handler, _ := setupEpicHandler(t)

	req := httptest.NewRequest("DELETE", "/invalid", nil)
	req = withChiURLParams(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEpicHandler_ProjectRoutes(t *testing.T) {
	handler, _ := setupEpicHandler(t)

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

func TestEpicHandler_Routes(t *testing.T) {
	handler, _ := setupEpicHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
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
