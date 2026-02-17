# Backend Phase 4: Wiki & Files

Rich text wiki system and file storage (local + S3).

## Task 4.1: Wiki Domain Models

**Commits:**
- `feat: add wiki page entity model`
- `feat: add wiki version history model`

**Package:** `internal/wiki/model.go`

## Task 4.2: Wiki Repository

**Commits:**
- `feat: add wiki repository interface`
- `feat: implement wiki CRUD operations`
- `feat: add wiki page versioning`
- `feat: add wiki tree/hierarchy queries`

**Package:** `internal/wiki/repository.go`

## Task 4.3: Wiki Service

**Commits:**
- `feat: add wiki service`
- `feat: add slug generation logic`
- `feat: add wiki page hierarchy management`

**Package:** `internal/wiki/service.go`

## Task 4.4: Wiki Handlers

**Commits:**
- `feat: add create wiki page handler`
- `feat: add list wiki pages handler`
- `feat: add get wiki page handler`
- `feat: add update wiki page handler`
- `feat: add delete wiki page handler`
- `feat: add wiki page version history handler`

**Package:** `internal/wiki/handler.go`

## Task 4.5: File Storage Abstraction

**Commits:**
- `feat: add storage interface`
- `feat: implement local filesystem storage`
- `feat: implement S3-compatible storage`
- `feat: add storage factory`

**Package:** `internal/file/storage.go`

## Task 4.6: File Entity Models

**Commits:**
- `feat: add file attachment entity model`

**Package:** `internal/file/model.go`

## Task 4.7: File Repository

**Commits:**
- `feat: add file repository interface`
- `feat: implement file metadata CRUD`

**Package:** `internal/file/repository.go`

## Task 4.8: File Service

**Commits:**
- `feat: add file service`
- `feat: add file upload logic with validation`
- `feat: add file download logic`

**Package:** `internal/file/service.go`

## Task 4.9: File Handlers

**Commits:**
- `feat: add file upload handler`
- `feat: add file download handler`
- `feat: add file delete handler`
- `feat: add list attachments handler`

**Package:** `internal/file/handler.go`

## Task 4.10: Database Migrations

**Commits:**
- `task: add wiki_pages table migration`
- `task: add wiki_page_versions table migration`
- `task: add file_attachments table migration`

---

## Total Backend Commits (Phase 4): 26
