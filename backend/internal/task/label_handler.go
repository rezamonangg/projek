package task

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type LabelHandler struct {
	taskService *TaskService
}

func NewLabelHandler(taskService *TaskService) *LabelHandler {
	return &LabelHandler{taskService: taskService}
}

func (h *LabelHandler) ProjectRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByProject)
	r.Post("/", h.Create)
	return r
}

func (h *LabelHandler) TaskRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListByTask)
	r.Post("/", h.AssignToTask)
	r.Delete("/{labelId}", h.RemoveFromTask)
	return r
}

func (h *LabelHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	labels, err := h.taskService.GetLabels(r.Context(), projectID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, labels)
}

func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_PROJECT_ID", "invalid project id")
		return
	}

	var req CreateLabelInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	label, err := h.taskService.CreateLabel(r.Context(), projectID, req)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusCreated, label)
}

func (h *LabelHandler) ListByTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "invalid task id")
		return
	}

	labels, err := h.taskService.GetTaskLabels(r.Context(), taskID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, labels)
}

type AssignLabelInput struct {
	LabelID uuid.UUID `json:"label_id" validate:"required"`
}

func (h *LabelHandler) AssignToTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "invalid task id")
		return
	}

	var req AssignLabelInput
	if err := common.ParseJSON(r, &req); err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := common.ValidateStruct(&req); err != nil {
		common.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.taskService.AddLabelToTask(r.Context(), taskID, req.LabelID); err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	task, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusOK, task)
}

func (h *LabelHandler) RemoveFromTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "invalid task id")
		return
	}

	labelIDStr := chi.URLParam(r, "labelId")
	labelID, err := uuid.Parse(labelIDStr)
	if err != nil {
		common.Error(w, http.StatusBadRequest, "INVALID_LABEL_ID", "invalid label id")
		return
	}

	if err := h.taskService.RemoveLabelFromTask(r.Context(), taskID, labelID); err != nil {
		common.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	common.Success(w, http.StatusNoContent, nil)
}
