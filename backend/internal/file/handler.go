package file

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Upload)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Get("/download", h.Download)
		r.Delete("/", h.Delete)
	})
	return r
}

func (h *Handler) ProjectRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByProject)
	return r
}

func (h *Handler) TaskRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByTask)
	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	common.Success(w, http.StatusOK, []FileAttachment{})
}

func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	attachments, err := h.service.GetByProject(r.Context(), projectID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, attachments)
}

func (h *Handler) ListByTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "invalid task id")
		return
	}

	attachments, err := h.service.GetByTask(r.Context(), taskID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, attachments)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "failed to parse multipart form")
		return
	}

	projectIDStr := r.FormValue("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_FILE", "file is required")
		return
	}
	defer file.Close()

	data := make([]byte, header.Size)
	if _, err := file.Read(data); err != nil {
		common.Error(w, http.StatusInternalServerError, "READ_ERROR", "failed to read file")
		return
	}

	var uploaderID uuid.UUID
	if member := getCurrentMember(r); member != nil {
		uploaderID = member.ID
	}

	input := CreateFileInput{
		ProjectID:   projectID,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		UploaderID:  uploaderID,
		Data:        data,
	}

	if taskIDStr := r.FormValue("task_id"); taskIDStr != "" {
		taskID, err := uuid.Parse(taskIDStr)
		if err == nil {
			input.TaskID = &taskID
		}
	}

	attachment, err := h.service.Upload(r.Context(), input)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, attachment)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid file id")
		return
	}

	url, err := h.service.GetURL(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "file not found")
		return
	}

	common.Success(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid file id")
		return
	}

	data, err := h.service.Download(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "file not found")
		return
	}

	w.Header().Set("Content-Disposition", "attachment")
	w.Write(data)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid file id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "file not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}

type Member struct {
	ID uuid.UUID `json:"id"`
}

func getCurrentMember(r *http.Request) *Member {
	if v := r.Context().Value("member"); v != nil {
		if m, ok := v.(*Member); ok {
			return m
		}
	}
	return nil
}
