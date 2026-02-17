package project

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testBoardID = uuid.MustParse("55555555-5555-5555-5555-555555555555")

func testBoard() *Board {
	now := time.Now()
	return &Board{
		ID:        testBoardID,
		ProjectID: testProjectID,
		Name:      "Main Board",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func setupBoardHandler(t *testing.T) (*BoardHandler, *MockBoardRepository) {
	ctrl := gomock.NewController(t)
	boardRepo := NewMockBoardRepository(ctrl)
	svc := NewService(nil, nil, boardRepo, nil, nil)
	return NewBoardHandler(svc), boardRepo
}

func TestBoardHandler_ListByProject_Success(t *testing.T) {
	handler, boardRepo := setupBoardHandler(t)

	boards := []Board{*testBoard()}
	boardRepo.EXPECT().GetByProject(gomock.Any(), testProjectID).Return(boards, nil)

	req := httptest.NewRequest("GET", "/"+testProjectID.String()+"/boards", nil)
	req = withChiURLParams(req, map[string]string{"projectId": testProjectID.String()})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp common.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestBoardHandler_ListByProject_InvalidProjectID(t *testing.T) {
	handler, _ := setupBoardHandler(t)

	req := httptest.NewRequest("GET", "/invalid/boards", nil)
	req = withChiURLParams(req, map[string]string{"projectId": "invalid"})

	rr := httptest.NewRecorder()
	handler.ListByProject(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBoardHandler_Create_Success(t *testing.T) {
	handler, boardRepo := setupBoardHandler(t)

	body := `{"name":"New Board"}`
	boardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/boards", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"projectId": testProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestBoardHandler_Create_MissingName(t *testing.T) {
	handler, _ := setupBoardHandler(t)

	body := `{}`

	req := httptest.NewRequest("POST", "/"+testProjectID.String()+"/boards", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"projectId": testProjectID.String()})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBoardHandler_Create_InvalidProjectID(t *testing.T) {
	handler, _ := setupBoardHandler(t)

	body := `{"name":"New Board"}`

	req := httptest.NewRequest("POST", "/invalid/boards", strings.NewReader(body))
	req = withChiURLParams(req, map[string]string{"projectId": "invalid"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBoardHandler_Get_Success(t *testing.T) {
	handler, boardRepo := setupBoardHandler(t)

	boardRepo.EXPECT().GetByID(gomock.Any(), testBoardID).Return(testBoard(), nil)

	req := httptest.NewRequest("GET", "/"+testBoardID.String(), nil)
	req = withChiURLParams(req, map[string]string{"id": testBoardID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestBoardHandler_Get_NotFound(t *testing.T) {
	handler, boardRepo := setupBoardHandler(t)

	boardRepo.EXPECT().GetByID(gomock.Any(), testBoardID).Return(nil, common.ErrNotFound)

	req := httptest.NewRequest("GET", "/"+testBoardID.String(), nil)
	req = withChiURLParams(req, map[string]string{"id": testBoardID.String()})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestBoardHandler_Get_InvalidID(t *testing.T) {
	handler, _ := setupBoardHandler(t)

	req := httptest.NewRequest("GET", "/invalid", nil)
	req = withChiURLParams(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBoardHandler_Delete_Success(t *testing.T) {
	handler, boardRepo := setupBoardHandler(t)

	boardRepo.EXPECT().Delete(gomock.Any(), testBoardID).Return(nil)

	req := httptest.NewRequest("DELETE", "/"+testBoardID.String(), nil)
	req = withChiURLParams(req, map[string]string{"id": testBoardID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestBoardHandler_Delete_NotFound(t *testing.T) {
	handler, boardRepo := setupBoardHandler(t)

	boardRepo.EXPECT().Delete(gomock.Any(), testBoardID).Return(common.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/"+testBoardID.String(), nil)
	req = withChiURLParams(req, map[string]string{"id": testBoardID.String()})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestBoardHandler_Delete_InvalidID(t *testing.T) {
	handler, _ := setupBoardHandler(t)

	req := httptest.NewRequest("DELETE", "/invalid", nil)
	req = withChiURLParams(req, map[string]string{"id": "invalid"})

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBoardHandler_ProjectRoutes(t *testing.T) {
	handler, _ := setupBoardHandler(t)

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

func TestBoardHandler_Routes(t *testing.T) {
	handler, _ := setupBoardHandler(t)

	router := handler.Routes()

	tests := []struct {
		method string
		path   string
	}{
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

func TestService_GetBoardByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	boardRepo := NewMockBoardRepository(ctrl)
	svc := NewService(nil, nil, boardRepo, nil, nil)

	expected := testBoard()
	boardRepo.EXPECT().GetByID(gomock.Any(), testBoardID).Return(expected, nil)

	board, err := svc.GetBoardByID(context.Background(), testBoardID)

	require.NoError(t, err)
	assert.Equal(t, expected, board)
}
