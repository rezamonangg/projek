# Phase 4: Wiki & File Attachments

Rich text wiki system and file storage (local + S3).

## Task 4.1: Wiki Domain Models

**Commits:**
- `feat: add wiki page entity model`
- `feat: add wiki version history model`

**Models:**
```go
type WikiPage struct {
    ID          uuid.UUID
    ProjectID   uuid.UUID
    Title       string
    Slug        string // URL-friendly
    Content     string // JSON from TipTap
    ParentID    *uuid.UUID // for nested pages
    CreatedBy   uuid.UUID
    UpdatedBy   uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type WikiPageVersion struct {
    ID        uuid.UUID
    PageID    uuid.UUID
    Content   string
    CreatedBy uuid.UUID
    CreatedAt time.Time
}
```

## Task 4.2: Wiki Repository

**Commits:**
- `feat: add wiki repository interface`
- `feat: implement wiki CRUD operations`
- `feat: add wiki page versioning`
- `feat: add wiki tree/hierarchy queries`

## Task 4.3: Wiki Service

**Commits:**
- `feat: add wiki service`
- `feat: add slug generation logic`
- `feat: add wiki page hierarchy management`

## Task 4.4: Wiki Handlers

**Commits:**
- `feat: add create wiki page handler`
- `feat: add list wiki pages handler (tree view)`
- `feat: add get wiki page handler`
- `feat: add update wiki page handler`
- `feat: add delete wiki page handler`
- `feat: add wiki page version history handler`

## Task 4.5: File Storage Abstraction

**Commits:**
- `feat: add storage interface`
- `feat: implement local filesystem storage`
- `feat: implement S3-compatible storage`
- `feat: add storage factory`

**Interface:**
```go
type Storage interface {
    Upload(ctx context.Context, key string, data []byte, contentType string) error
    Download(ctx context.Context, key string) ([]byte, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) (string, error)
}
```

## Task 4.6: File Entity Models

**Commits:**
- `feat: add file attachment entity model`
- `feat: add file metadata model`

**Models:**
```go
type FileAttachment struct {
    ID          uuid.UUID
    CommunityID uuid.UUID
    ProjectID   *uuid.UUID // nullable
    TaskID      *uuid.UUID // nullable
    WikiPageID  *uuid.UUID // nullable
    
    FileName    string
    FileSize    int64
    ContentType string
    StorageKey  string
    StorageType string // local or s3
    
    UploadedBy  uuid.UUID
    CreatedAt   time.Time
}
```

## Task 4.7: File Repository

**Commits:**
- `feat: add file repository interface`
- `feat: implement file metadata CRUD`
- `feat: add file listing by entity`

## Task 4.8: File Service

**Commits:**
- `feat: add file service`
- `feat: add file upload logic with size validation`
- `feat: add file type validation`
- `feat: add file download logic`

**Validation:**
- Max file size: configurable (default 10MB)
- Allowed types: images, documents, common formats

## Task 4.9: File Handlers

**Commits:**
- `feat: add file upload handler`
- `feat: add file download handler`
- `feat: add file delete handler`
- `feat: add list attachments handler`

## Task 4.10: Task Attachments

**Commits:**
- `feat: add attach file to task endpoint`
- `feat: add remove attachment from task endpoint`
- `feat: add list task attachments endpoint`

## Task 4.11: Wiki Attachments

**Commits:**
- `feat: add attach file to wiki endpoint`
- `feat: add remove attachment from wiki endpoint`
- `feat: add list wiki attachments endpoint`

## Task 4.12: Database Migrations for Wiki & Files

**Commits:**
- `task: add wiki_pages table migration`
- `task: add wiki_page_versions table migration`
- `task: add file_attachments table migration`
- `task: add indexes for file lookups`

---

## Total Commits: 26
**Estimated Time:** 3-4 days
