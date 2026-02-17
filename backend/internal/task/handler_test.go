package task

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

func setupTaskHandler(t *testing.T) (*Handler, *MockRepository) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	svc := NewService(repo, nil)
	return NewHandler(svc), repo
}

func withTaskChiURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

func TestTaskHandler_ListByBoard_Success(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	tasks := []Task{*testTask()}
	repo.EXPECT().GetByBoard(gomock.Any(), testBoardID, TaskFilter{}).Return(tasks, nil)

	req := httptest.NewRequest("GET", "/"+testBoardID.String()+"/tasks", nil)
	req = withTaskChiURLParams(req, map[string]string{"boardId": testBoardID.String()})

	rr := httptest.NewRecorder()
	handler.ListByBoard(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestTaskHandler_ListByBoard_InvalidBoardID(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	req := httptest.NewRequest("GET", "/invalid/tasks", nil)
	req = withTaskChiURLParams(req, map[string]string{"boardId": "invalid"})

	rr := httptest.NewRecorder()
	handler.ListByBoard(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTaskHandler_Create_Success(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	body := `{"title":"New Task","description":"Desc","reporter_id":"33333333-3333-3333-3333-333333333333"}`
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest("POST", "/"+testBoardID.String()+"/tasks", strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"boardId": testBoardID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestTaskHandler_Create_MissingTitle(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	body := `{"description":"Desc","reporter_id":"33333333-3333-3333-3333-333333333333"}`

	req := httptest.NewRequest("POST", "/"+testBoardID.String()+"/tasks", strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"boardId": testBoardID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTaskHandler_Create_InvalidBoardID(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	body := `{"title":"New Task","reporter_id":"33333333-3333-3333-3333-333333333333"}`

	req := httptest.NewRequest("POST", "/invalid/tasks", strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"boardId": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTaskHandler_Get_Success(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(testTask(), nil)

	req := httptest.NewRequest("GET", "/"+testTaskID.String(), nil)
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTaskHandler_Get_NotFound(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(nil, common.ErrNotFound)

	req := httptest.NewRequest("GET", "/"+testTaskID.String(), nil)
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestTaskHandler_Get_InvalidID(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	req := httptest.NewRequest("GET", "/invalid", nil)
	req = withTaskChiURLParams(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTaskHandler_Update_Success(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	existing := testTask()
	repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"title":"Updated Task","status":"inprogress"}`
	req := httptest.NewRequest("PUT", "/"+testTaskID.String(), strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTaskHandler_Update_NotFound(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(nil, common.ErrNotFound)

	body := `{"title":"Updated Task"}`
	req := httptest.NewRequest("PUT", "/"+testTaskID.String(), strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestTaskHandler_Move_Success(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	existing := testTask()
	repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body := `{"status":"done","position":5}`
	req := httptest.NewRequest("PATCH", "/"+testTaskID.String()+"/move", strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Move(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTaskHandler_Move_InvalidStatus(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	body := `{"status":"invalid","position":5}`
	req := httptest.NewRequest("PATCH", "/"+testTaskID.String()+"/move", strings.NewReader(body))
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Move(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTaskHandler_Delete_Success(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testTaskID).Return(nil)

	req := httptest.NewRequest("DELETE", "/"+testTaskID.String(), nil)
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestTaskHandler_Delete_NotFound(t *testing.T) {
	handler, repo := setupTaskHandler(t)

	repo.EXPECT().Delete(gomock.Any(), testTaskID).Return(common.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/"+testTaskID.String(), nil)
	req = withTaskChiURLParams(req, map[string]string{"id": testTaskID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestTaskHandler_BoardRoutes(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	router := handler.BoardRoutes()

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

func TestTaskHandler_Routes(t *testing.T) {
	handler, _ := setupTaskHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/{id}"},
		{"PUT", "/{id}"},
		{"DELETE", "/{id}"},
		{"PATCH", "/{id}/move"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var body *strings.Reader
			if tt.method == "PUT" || tt.method == "PATCH" {
				body = strings.NewReader(`{}`)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.method == "PUT" || tt.method == "PATCH" {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.NotEqual(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	svc := NewService(repo, nil)

	input := CreateTaskInput{
		Title:      "New Task",
		ReporterID: testMemberID,
		Status:     StatusBacklog,
	}

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	task, err := svc.Create(context.Background(), testBoardID, input)

	require.NoError(t, err)
	assert.Equal(t, "New Task", task.Title)
	assert.Equal(t, testBoardID, task.BoardID)
	assert.Equal(t, StatusBacklog, task.Status)
}

func TestService_GetByBoard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	svc := NewService(repo, nil)

	tasks := []Task{*testTask()}
	repo.EXPECT().GetByBoard(gomock.Any(), testBoardID, TaskFilter{}).Return(tasks, nil)

	result, err := svc.GetByBoard(context.Background(), testBoardID, TaskFilter{})

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestService_Move_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	svc := NewService(repo, nil)

	existing := testTask()
	repo.EXPECT().GetByID(gomock.Any(), testTaskID).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	task, err := svc.Move(context.Background(), testTaskID, StatusDone, 10)

	require.NoError(t, err)
	assert.Equal(t, StatusDone, task.Status)
	assert.Equal(t, 10, task.Position)
}
