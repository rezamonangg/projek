# Backend Implementation Plan: Files Module

**Priority:** LOW
**Estimated Time:** 2-3 hours
**Dependencies:** None (but relates to tasks for attachments)

## Current State

### What Exists
- Nothing in backend

### What's Missing
- Model definitions
- Repository + interface
- Service
- Handler
- Storage implementation (local filesystem or cloud)
- Database migration

## Frontend API Contract

Based on `frontend/src/lib/api/real/files.ts`:

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| POST | `/files/upload` | FormData: file, uploader_id, task_id? | `Attachment` |
| GET | `/files/{id}` | - | File blob |
| GET | `/files?task_id=...` | - | `Attachment[]` |

## Implementation Steps

### Step 1: Create Module Directory

```bash
mkdir -p backend/internal/file
```

### Step 2: Define Model

**File:** `backend/internal/file/model.go`

```go
package file

import (
    "time"
    "github.com/google/uuid"
)

type Attachment struct {
    ID              uuid.UUID `json:"id"`
    ProjectID       uuid.UUID `json:"project_id"`
    TaskID          *uuid.UUID `json:"task_id,omitempty"`
    UploaderID      uuid.UUID `json:"uploader_id"`
    Filename        string    `json:"filename"`
    OriginalFilename string   `json:"original_filename"`
    ContentType     string    `json:"content_type"`
    Size            int64     `json:"size"`
    StoragePath     string    `json:"storage_path"`
    StorageType     string    `json:"storage_type"` // "local", "s3", etc.
    CreatedAt       time.Time `json:"created_at"`
}
```

### Step 3: Create Storage Interface

**File:** `backend/internal/file/storage.go`

```go
package file

type Storage interface {
    Save(ctx context.Context, filename string, data []byte) (path string, err error)
    Get(ctx context.Context, path string) ([]byte, error)
    Delete(ctx context.Context, path string) error
}

type LocalStorage struct {
    basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
    return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) Save(ctx context.Context, filename string, data []byte) (string, error) {
    // Generate unique filename, save to disk
}

func (s *LocalStorage) Get(ctx context.Context, path string) ([]byte, error) {
    // Read from disk
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
    // Delete from disk
}
```

### Step 4: Create Repository

**File:** `backend/internal/file/repository.go`

```go
package file

type Repository interface {
    Create(ctx context.Context, attachment *Attachment) error
    GetByID(ctx context.Context, id uuid.UUID) (*Attachment, error)
    GetByTask(ctx context.Context, taskID uuid.UUID) ([]Attachment, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Step 5: Create Service

**File:** `backend/internal/file/service.go`

```go
package file

import (
    "io"
    "mime/multipart"
)

type Service struct {
    repo    Repository
    storage Storage
}

func NewService(repo Repository, storage Storage) *Service {
    return &Service{repo: repo, storage: storage}
}

func (s *Service) Upload(ctx context.Context, projectID uuid.UUID, taskID *uuid.UUID, uploaderID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*Attachment, error) {
    // Read file content
    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }
    
    // Generate storage path
    storagePath, err := s.storage.Save(ctx, header.Filename, data)
    if err != nil {
        return nil, err
    }
    
    // Create attachment record
    attachment := &Attachment{
        ID:               uuid.New(),
        ProjectID:        projectID,
        TaskID:           taskID,
        UploaderID:       uploaderID,
        Filename:         generateUniqueFilename(header.Filename),
        OriginalFilename: header.Filename,
        ContentType:      header.Header.Get("Content-Type"),
        Size:             header.Size,
        StoragePath:      storagePath,
        StorageType:      "local",
        CreatedAt:        common.Now(),
    }
    
    if err := s.repo.Create(ctx, attachment); err != nil {
        s.storage.Delete(ctx, storagePath) // Cleanup
        return nil, err
    }
    
    return attachment, nil
}

func (s *Service) Download(ctx context.Context, id uuid.UUID) ([]byte, *Attachment, error) {
    attachment, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, nil, err
    }
    
    data, err := s.storage.Get(ctx, attachment.StoragePath)
    if err != nil {
        return nil, nil, err
    }
    
    return data, attachment, nil
}

func (s *Service) GetByTask(ctx context.Context, taskID uuid.UUID) ([]Attachment, error) {
    return s.repo.GetByTask(ctx, taskID)
}
```

### Step 6: Create Handler

**File:** `backend/internal/file/handler.go`

```go
package file

import (
    "net/http"
    "strconv"
)

const maxUploadSize = 10 * 1024 * 1024 // 10MB

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Post("/upload", h.Upload)
    r.Get("/{id}", h.Download)
    r.Get("/", h.List)
    return r
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
    // Limit request size
    r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
    
    // Parse multipart form
    if err := r.ParseMultipartForm(maxUploadSize); err != nil {
        common.Error(w, 400, "INVALID_REQUEST", "file too large")
        return
    }
    
    // Get form values
    uploaderID := r.FormValue("uploader_id")
    taskID := r.FormValue("task_id")
    
    // Get file
    file, header, err := r.FormFile("file")
    if err != nil {
        common.Error(w, 400, "INVALID_REQUEST", "no file provided")
        return
    }
    defer file.Close()
    
    // Upload
    var taskUUID *uuid.UUID
    if taskID != "" {
        id := uuid.MustParse(taskID)
        taskUUID = &id
    }
    
    // Get project_id from context or form value
    projectID := uuid.MustParse(r.FormValue("project_id"))
    
    attachment, err := h.service.Upload(r.Context(), projectID, taskUUID, uuid.MustParse(uploaderID), file, header)
    if err != nil {
        common.Error(w, 500, "INTERNAL_ERROR", err.Error())
        return
    }
    
    common.Success(w, 201, attachment)
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    data, attachment, err := h.service.Download(r.Context(), uuid.MustParse(id))
    if err != nil {
        common.Error(w, 404, "NOT_FOUND", "file not found")
        return
    }
    
    w.Header().Set("Content-Type", attachment.ContentType)
    w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, attachment.OriginalFilename))
    w.Header().Set("Content-Length", strconv.FormatInt(attachment.Size, 10))
    w.Write(data)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    taskID := r.URL.Query().Get("task_id")
    
    if taskID != "" {
        attachments, err := h.service.GetByTask(r.Context(), uuid.MustParse(taskID))
        if err != nil {
            common.Error(w, 500, "INTERNAL_ERROR", err.Error())
            return
        }
        common.Success(w, 200, attachments)
        return
    }
    
    common.Success(w, 200, []Attachment{})
}
```

### Step 7: Create Migration

**File:** `backend/db/migrations/YYYYMMDDHHMMSS_create_attachments.sql`

```sql
-- migrate:up
CREATE TABLE attachments (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    uploader_id UUID NOT NULL REFERENCES members(id),
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(255),
    size BIGINT NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    storage_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_attachments_task_id ON attachments(task_id);
CREATE INDEX idx_attachments_project_id ON attachments(project_id);

-- migrate:down
DROP TABLE IF EXISTS attachments;
```

### Step 8: Register Routes

**File:** `backend/internal/router/router.go`

```go
// Configure storage
storage := file.NewLocalStorage("/var/lib/projek/uploads") // or from config

fileRepo := file.NewRepository(db)
fileService := file.NewService(fileRepo, storage)
fileHandler := file.NewHandler(fileService)
r.Mount("/files", fileHandler.Routes())
```

### Step 9: Create Upload Directory

```bash
mkdir -p /var/lib/projek/uploads
# Or configure in .env and create at runtime
```

## Configuration

Add to config:
```go
type Config struct {
    // ...
    UploadPath string `env:"UPLOAD_PATH" envDefault:"./uploads"`
}
```

## Response Format

```json
{
    "id": "uuid",
    "project_id": "uuid",
    "task_id": "uuid or null",
    "uploader_id": "uuid",
    "filename": "abc123.pdf",
    "original_filename": "document.pdf",
    "content_type": "application/pdf",
    "size": 1024,
    "storage_path": "/uploads/2024/01/abc123.pdf",
    "storage_type": "local",
    "created_at": 1704067200
}
```

## Unit Tests

### Test File Structure

**File:** `backend/internal/file/handler_test.go`

```go
package file

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/monachy/projek/internal/common"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

var testAttachmentID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
var testUploaderID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

func testAttachment() *Attachment {
    return &Attachment{
        ID:               testAttachmentID,
        ProjectID:        testProjectID,
        UploaderID:       testUploaderID,
        Filename:         "test-123.txt",
        OriginalFilename: "test.txt",
        ContentType:      "text/plain",
        Size:             100,
        StoragePath:      "/uploads/2024/01/test-123.txt",
        StorageType:      "local",
        CreatedAt:        time.Now(),
    }
}

func setupFileHandler(t *testing.T) (*Handler, *MockRepository, *MockStorage) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    storage := NewMockStorage(ctrl)
    svc := NewService(repo, storage)
    return NewHandler(svc), repo, storage
}

func createMultipartBody(t *testing.T, filename, content string) (*bytes.Buffer, string) {
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)
    
    part, err := writer.CreateFormFile("file", filename)
    require.NoError(t, err)
    _, err = io.WriteString(part, content)
    require.NoError(t, err)
    
    _ = writer.WriteField("uploader_id", testUploaderID.String())
    _ = writer.WriteField("project_id", testProjectID.String())
    
    err = writer.Close()
    require.NoError(t, err)
    
    return &buf, writer.FormDataContentType()
}
```

### Handler Tests

```go
func TestFileHandler_Upload_Success(t *testing.T) {
    handler, repo, storage := setupFileHandler(t)
    
    buf, contentType := createMultipartBody(t, "test.txt", "test content")
    
    storage.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return("/uploads/test.txt", nil)
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    req := httptest.NewRequest("POST", "/files/upload", buf)
    req.Header.Set("Content-Type", contentType)
    
    rr := httptest.NewRecorder()
    handler.Upload(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
    
    var resp Attachment
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Equal(t, "test.txt", resp.OriginalFilename)
}

func TestFileHandler_Upload_NoFile(t *testing.T) {
    handler, _, _ := setupFileHandler(t)
    
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)
    _ = writer.WriteField("uploader_id", testUploaderID.String())
    writer.Close()
    
    req := httptest.NewRequest("POST", "/files/upload", &buf)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    rr := httptest.NewRecorder()
    handler.Upload(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFileHandler_Upload_TooLarge(t *testing.T) {
    handler, _, _ := setupFileHandler(t)
    
    // Create a body larger than maxUploadSize
    largeContent := make([]byte, maxUploadSize+1)
    buf, contentType := createMultipartBody(t, "large.txt", string(largeContent))
    
    req := httptest.NewRequest("POST", "/files/upload", buf)
    req.Header.Set("Content-Type", contentType)
    
    rr := httptest.NewRecorder()
    handler.Upload(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFileHandler_Download_Success(t *testing.T) {
    handler, repo, storage := setupFileHandler(t)
    
    attachment := testAttachment()
    fileContent := []byte("test file content")
    
    repo.EXPECT().GetByID(gomock.Any(), testAttachmentID).Return(attachment, nil)
    storage.EXPECT().Get(gomock.Any(), attachment.StoragePath).Return(fileContent, nil)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Download)
    
    req := httptest.NewRequest("GET", "/"+testAttachmentID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
    assert.Contains(t, rr.Header().Get("Content-Disposition"), "test.txt")
    assert.Equal(t, fileContent, rr.Body.Bytes())
}

func TestFileHandler_Download_NotFound(t *testing.T) {
    handler, repo, _ := setupFileHandler(t)
    
    repo.EXPECT().GetByID(gomock.Any(), testAttachmentID).Return(nil, common.ErrNotFound)
    
    r := chi.NewRouter()
    r.Get("/{id}", handler.Download)
    
    req := httptest.NewRequest("GET", "/"+testAttachmentID.String(), nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestFileHandler_List_ByTask(t *testing.T) {
    handler, repo, _ := setupFileHandler(t)
    
    attachments := []Attachment{*testAttachment()}
    repo.EXPECT().GetByTask(gomock.Any(), testTaskID).Return(attachments, nil)
    
    req := httptest.NewRequest("GET", "/files?task_id="+testTaskID.String(), nil)
    
    rr := httptest.NewRecorder()
    handler.List(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []Attachment
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
}

func TestFileHandler_List_Empty(t *testing.T) {
    handler, _, _ := setupFileHandler(t)
    
    req := httptest.NewRequest("GET", "/files", nil)
    
    rr := httptest.NewRecorder()
    handler.List(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp []Attachment
    json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.Empty(t, resp)
}
```

### Service Tests

**File:** `backend/internal/file/service_test.go`

```go
func TestService_Upload_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    storage := NewMockStorage(ctrl)
    svc := NewService(repo, storage)
    
    fileContent := []byte("test content")
    header := &multipart.FileHeader{
        Filename: "test.txt",
        Size:     int64(len(fileContent)),
        Header:   make(textproto.MIMEHeader),
    }
    header.Header.Set("Content-Type", "text/plain")
    
    storage.EXPECT().Save(gomock.Any(), gomock.Any(), fileContent).Return("/uploads/test.txt", nil)
    repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    attachment, err := svc.Upload(
        context.Background(),
        testProjectID,
        nil,
        testUploaderID,
        io.NopCloser(bytes.NewReader(fileContent)),
        header,
    )
    
    require.NoError(t, err)
    assert.Equal(t, "test.txt", attachment.OriginalFilename)
    assert.Equal(t, testProjectID, attachment.ProjectID)
}

func TestService_Upload_StorageFails(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    storage := NewMockStorage(ctrl)
    svc := NewService(repo, storage)
    
    fileContent := []byte("test content")
    header := &multipart.FileHeader{Filename: "test.txt"}
    
    storage.EXPECT().Save(gomock.Any(), gomock.Any(), fileContent).Return("", errors.New("disk full"))
    
    attachment, err := svc.Upload(
        context.Background(),
        testProjectID,
        nil,
        testUploaderID,
        io.NopCloser(bytes.NewReader(fileContent)),
        header,
    )
    
    assert.Nil(t, attachment)
    assert.Error(t, err)
}

func TestService_Download_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    repo := NewMockRepository(ctrl)
    storage := NewMockStorage(ctrl)
    svc := NewService(repo, storage)
    
    attachment := testAttachment()
    fileContent := []byte("test content")
    
    repo.EXPECT().GetByID(gomock.Any(), testAttachmentID).Return(attachment, nil)
    storage.EXPECT().Get(gomock.Any(), attachment.StoragePath).Return(fileContent, nil)
    
    data, att, err := svc.Download(context.Background(), testAttachmentID)
    
    require.NoError(t, err)
    assert.Equal(t, fileContent, data)
    assert.Equal(t, attachment, att)
}
```

### Storage Tests

**File:** `backend/internal/file/storage_test.go`

```go
func TestLocalStorage_Save_Success(t *testing.T) {
    tmpDir := t.TempDir()
    storage := NewLocalStorage(tmpDir)
    
    content := []byte("test content")
    path, err := storage.Save(context.Background(), "test.txt", content)
    
    require.NoError(t, err)
    assert.Contains(t, path, "test")
    
    // Verify file exists
    _, err = os.Stat(filepath.Join(tmpDir, filepath.Base(path)))
    require.NoError(t, err)
}

func TestLocalStorage_Get_Success(t *testing.T) {
    tmpDir := t.TempDir()
    storage := NewLocalStorage(tmpDir)
    
    // First save a file
    content := []byte("test content")
    path, _ := storage.Save(context.Background(), "test.txt", content)
    
    // Then get it
    data, err := storage.Get(context.Background(), path)
    
    require.NoError(t, err)
    assert.Equal(t, content, data)
}

func TestLocalStorage_Delete_Success(t *testing.T) {
    tmpDir := t.TempDir()
    storage := NewLocalStorage(tmpDir)
    
    content := []byte("test content")
    path, _ := storage.Save(context.Background(), "test.txt", content)
    
    err := storage.Delete(context.Background(), path)
    
    require.NoError(t, err)
    
    // Verify file is gone
    _, err = os.Stat(filepath.Join(tmpDir, filepath.Base(path)))
    assert.True(t, os.IsNotExist(err))
}

func TestLocalStorage_Get_NotFound(t *testing.T) {
    tmpDir := t.TempDir()
    storage := NewLocalStorage(tmpDir)
    
    _, err := storage.Get(context.Background(), "/nonexistent/file.txt")
    
    assert.Error(t, err)
}
```

### Test Cases Summary

| Test | Description | Expected Status |
|------|-------------|-----------------|
| `TestFileHandler_Upload_Success` | Uploads file | 201 |
| `TestFileHandler_Upload_NoFile` | No file in request | 400 |
| `TestFileHandler_Upload_TooLarge` | File exceeds limit | 400 |
| `TestFileHandler_Download_Success` | Downloads file with headers | 200 |
| `TestFileHandler_Download_NotFound` | File doesn't exist | 404 |
| `TestFileHandler_List_ByTask` | Lists by task ID | 200 |
| `TestFileHandler_List_Empty` | No task ID param | 200 |

## Files to Create

| File | Action |
|------|--------|
| `backend/internal/file/model.go` | CREATE |
| `backend/internal/file/storage.go` | CREATE |
| `backend/internal/file/storage_test.go` | CREATE |
| `backend/internal/file/repository.go` | CREATE |
| `backend/internal/file/repository_test.go` | CREATE |
| `backend/internal/file/service.go` | CREATE |
| `backend/internal/file/service_test.go` | CREATE |
| `backend/internal/file/handler.go` | CREATE |
| `backend/internal/file/handler_test.go` | CREATE |
| `backend/db/migrations/*_create_attachments.sql` | CREATE |
| `backend/internal/router/router.go` | MODIFY |

## Security Considerations

1. **File size limit** - 10MB max
2. **File type validation** - Consider restricting dangerous types
3. **Path traversal** - Sanitize filenames
4. **Access control** - Verify user has access to project/task
5. **Virus scanning** - Consider for production

## Verification

```bash
cd backend && go test ./internal/file/... -v
```
