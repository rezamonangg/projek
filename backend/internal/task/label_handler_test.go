package task

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

var testLabelID = uuid.MustParse("88888888-8888-8888-8888-888888888888")
var testLabelProjectID = uuid.MustParse("99999999-9999-9999-9999-999999999999")

func testLabel() *Label {
	return &Label{
		ID:        testLabelID,
		ProjectID: testLabelProjectID,
		Name:      "Bug",
		Color:     "#FF0000",
		CreatedAt: time.Now(),
	}
}

func setupLabelHandler(t *testing.T) (*LabelHandler, *MockLabelRepository, *MockRepository) {
	ctrl := gomock.NewController(t)
	labelRepo := NewMockLabelRepository(ctrl)
	taskRepo := NewMockRepository(ctrl)
	svc := NewTaskService(taskRepo, labelRepo)
	return NewLabelHandler(svc), labelRepo, taskRepo
}

func withLabelChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestLabelHandler_ListByProject_Success(t *testing.T) {
	handler, labelRepo, _ := setupLabelHandler(t)

	labels := []Label{*testLabel()}
	labelRepo.EXPECT().GetByProject(gomock.Any(), testLabelProjectID).Return(labels, nil)

	req := httptest.NewRequest("GET", "/"+testLabelProjectID.String()+"/labels", nil)
	req = withLabelChiURLParams(req, map[string]string{"projectId": testLabelProjectID.String()})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestLabelHandler_ListByProject_InvalidProjectID(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	req := httptest.NewRequest("GET", "/invalid/labels", nil)
	req = withLabelChiURLParams(req, map[string]string{"projectId": "invalid"})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLabelHandler_Create_Success(t *testing.T) {
	handler, labelRepo, _ := setupLabelHandler(t)

	body := `{"name":"Feature","color":"#00FF00"}`
	labelRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest("POST", "/"+testLabelProjectID.String()+"/labels", strings.NewReader(body))
	req = withLabelChiURLParams(req, map[string]string{"projectId": testLabelProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestLabelHandler_Create_MissingFields(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	tests := []struct {
		name string
		body string
	}{
		{"missing name", `{"color":"#00FF00"}`},
		{"empty body", `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/"+testLabelProjectID.String()+"/labels", strings.NewReader(tt.body))
			req = withLabelChiURLParams(req, map[string]string{"projectId": testLabelProjectID.String()})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Create(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestLabelHandler_Create_InvalidProjectID(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	body := `{"name":"Feature","color":"#00FF00"}`

	req := httptest.NewRequest("POST", "/invalid/labels", strings.NewReader(body))
	req = withLabelChiURLParams(req, map[string]string{"projectId": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLabelHandler_ListByTask_Success(t *testing.T) {
	handler, labelRepo, _ := setupLabelHandler(t)

	labels := []Label{*testLabel()}
	labelRepo.EXPECT().GetByTask(gomock.Any(), testTaskID).Return(labels, nil)

	req := httptest.NewRequest("GET", "/"+testTaskID.String()+"/labels", nil)
	req = withLabelChiURLParams(req, map[string]string{"taskId": testTaskID.String()})

	rr := httptest.NewRecorder()
	handler.ListByTask(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLabelHandler_AssignToTask_Success(t *testing.T) {
	handler, labelRepo, taskRepo := setupLabelHandler(t)

	labelRepo.EXPECT().AddToTask(gomock.Any(), testTaskID, testLabelID).Return(nil)
	taskRepo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(testTask(), nil)

	body := `{"label_id":"88888888-8888-8888-8888-888888888888"}`
	req := httptest.NewRequest("POST", "/"+testTaskID.String()+"/labels", strings.NewReader(body))
	req = withLabelChiURLParams(req, map[string]string{"taskId": testTaskID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.AssignToTask(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLabelHandler_AssignToTask_InvalidTaskID(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	body := `{"label_id":"88888888-8888-8888-8888-888888888888"}`
	req := httptest.NewRequest("POST", "/invalid/labels", strings.NewReader(body))
	req = withLabelChiURLParams(req, map[string]string{"taskId": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.AssignToTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLabelHandler_RemoveFromTask_Success(t *testing.T) {
	handler, labelRepo, _ := setupLabelHandler(t)

	labelRepo.EXPECT().RemoveFromTask(gomock.Any(), testTaskID, testLabelID).Return(nil)

	req := httptest.NewRequest("DELETE", "/"+testTaskID.String()+"/labels/"+testLabelID.String(), nil)
	req = withLabelChiURLParams(req, map[string]string{"taskId": testTaskID.String(), "labelId": testLabelID.String()})

	rr := httptest.NewRecorder()
	handler.RemoveFromTask(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestLabelHandler_RemoveFromTask_InvalidTaskID(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	req := httptest.NewRequest("DELETE", "/invalid/labels/"+testLabelID.String(), nil)
	req = withLabelChiURLParams(req, map[string]string{"taskId": "invalid", "labelId": testLabelID.String()})

	rr := httptest.NewRecorder()
	handler.RemoveFromTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLabelHandler_RemoveFromTask_InvalidLabelID(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	req := httptest.NewRequest("DELETE", "/"+testTaskID.String()+"/labels/invalid", nil)
	req = withLabelChiURLParams(req, map[string]string{"taskId": testTaskID.String(), "labelId": "invalid"})

	rr := httptest.NewRecorder()
	handler.RemoveFromTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLabelHandler_ProjectRoutes(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

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

func TestLabelHandler_TaskRoutes(t *testing.T) {
	handler, _, _ := setupLabelHandler(t)

	router := handler.TaskRoutes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"POST", "/"},
		{"DELETE", "/{labelId}"},
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
