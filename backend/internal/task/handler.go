package task

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

func (h *Handler) BoardRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByBoard)
	r.Post("/", h.Create)
	return r
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
		r.Patch("/move", h.Move)
	})
	return r
}

type MoveTaskInput struct {
	Status   TaskStatus `json:"status"`
	Position int        `json:"position"`
}

func (h *Handler) ListByBoard(w http.ResponseWriter, r *http.Request) {
	boardIDStr := chi.URLParam(r, "boardId")
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_BOARD_ID", "invalid board id")
		return
	}

	filter := TaskFilter{}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = TaskStatus(status)
	}
	if assigneeID := r.URL.Query().Get("assignee_id"); assigneeID != "" {
		if id, err := uuid.Parse(assigneeID); err == nil {
			filter.AssigneeID = id
		}
	}
	if epicID := r.URL.Query().Get("epic_id"); epicID != "" {
		if id, err := uuid.Parse(epicID); err == nil {
			filter.EpicID = id
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = search
	}

	tasks, err := h.service.GetByBoard(r.Context(), boardID, filter)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, tasks)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	boardIDStr := chi.URLParam(r, "boardId")
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_BOARD_ID", "invalid board id")
		return
	}

	var req CreateTaskInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	task, err := h.service.Create(r.Context(), boardID, req)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, task)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid task id")
		return
	}

	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}

	common.Success(w, http.StatusOK, task)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid task id")
		return
	}

	var req UpdateTaskInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	task, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}

	common.Success(w, http.StatusOK, task)
}

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid task id")
		return
	}

	var req MoveTaskInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if !req.Status.IsValid() {
		common.Error(w, http.StatusBadRequest, "INVALID_STATUS", "invalid task status")
		return
	}

	task, err := h.service.Move(r.Context(), id, req.Status, req.Position)
	if err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}

	common.Success(w, http.StatusOK, task)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid task id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		common.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}
