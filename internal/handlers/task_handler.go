package handlers

import (
	"net/http"
	"strconv"
	"task-backend/internal/dto"
	"task-backend/internal/res"
	"task-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskService *services.TaskService
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			Message: "Failed to retrieve user ID",
			Status:  http.StatusInternalServerError,
			Error:   "invalid id",
		})
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			Message: "Invalid user ID type",
			Status:  http.StatusInternalServerError,
			Error:   "type assertion failed",
		})
		return
	}

	tasks := h.TaskService.GetAllTasks(ctx, userID)
	c.JSON(http.StatusOK, res.SuccessResponse{
		Message: "All tasks retrieved",
		Status:  http.StatusOK,
		Data:    tasks,
	})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			Message: "Failed to retrieve user ID",
			Status:  http.StatusInternalServerError,
			Error:   "invalid id",
		})
		return
	}
	userID := userIDRaw.(string)

	taskID := c.Param("id")
	taskIDUint, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			Message: "Invalid task ID",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	task, found := h.TaskService.GetTaskByID(c, userID, taskIDUint)
	if !found {
		c.JSON(http.StatusNotFound, res.ErrorResponse{
			Message: "Task not found",
			Status:  http.StatusNotFound,
			Error:   "invalid id",
		})
		return
	}

	c.JSON(http.StatusOK, res.SuccessResponse{
		Message: "Task retrieved",
		Status:  http.StatusOK,
		Data:    task,
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var newTask dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			Message: "Invalid request body",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	if err := newTask.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			Message: "Validation failed",
			Status:  http.StatusBadRequest,
			Error:   err,
		})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			Message: "Failed to retrieve user ID",
			Status:  http.StatusInternalServerError,
			Error:   "invalid id",
		})
		return
	}
	userID := userIDRaw.(string)

	created := h.TaskService.CreateTask(c, userID, newTask)

	c.JSON(http.StatusCreated, res.SuccessResponse{
		Message: "Task created",
		Status:  http.StatusCreated,
		Data:    created,
	})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			Message: "Failed to retrieve user ID",
			Status:  http.StatusInternalServerError,
			Error:   "invalid id",
		})
		return
	}
	userID := userIDRaw.(string)

	taskID := c.Param("id")
	taskIDUint, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			Message: "Invalid task ID",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	var updateData dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			Message: "Invalid request body",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	updated, ok := h.TaskService.UpdateTask(c, userID, taskIDUint, updateData)
	if !ok {
		c.JSON(http.StatusNotFound, res.ErrorResponse{
			Message: "Task not found",
			Status:  http.StatusNotFound,
			Error:   "invalid id",
		})
		return
	}

	c.JSON(http.StatusOK, res.SuccessResponse{
		Message: "Task updated",
		Status:  http.StatusOK,
		Data:    updated,
	})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, res.ErrorResponse{
			Message: "Failed to retrieve user ID",
			Status:  http.StatusInternalServerError,
			Error:   "invalid id",
		})
		return
	}
	userID := userIDRaw.(string)

	taskID := c.Param("id")
	taskIDUint, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.ErrorResponse{
			Message: "Invalid task ID",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	ok := h.TaskService.DeleteTask(c, userID, taskIDUint)
	if !ok {
		c.JSON(http.StatusNotFound, res.ErrorResponse{
			Message: "Task not found",
			Status:  http.StatusNotFound,
			Error:   "invalid id",
		})
		return
	}

	c.JSON(http.StatusOK, res.SuccessResponse{
		Message: "Task deleted",
		Status:  http.StatusOK,
		Data:    nil,
	})
}
