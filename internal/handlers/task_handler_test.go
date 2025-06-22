package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-backend/internal/dto"
	"task-backend/internal/handlers"
	"task-backend/internal/models"
	"task-backend/internal/res"
	"task-backend/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTaskRepository struct {
	tasks  map[string]map[uint64]models.Task
	nextID uint64
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks:  make(map[string]map[uint64]models.Task),
		nextID: 1,
	}
}

func (m *MockTaskRepository) Create(userID string, task models.Task) models.Task {
	if m.tasks[userID] == nil {
		m.tasks[userID] = make(map[uint64]models.Task)
	}
	task.ID = m.nextID
	m.nextID++
	task.UserID = userID
	m.tasks[userID][task.ID] = task
	return task
}

func (m *MockTaskRepository) GetAll(userID string) []models.Task {
	tasksMap := m.tasks[userID]
	tasks := make([]models.Task, 0, len(tasksMap))
	for _, t := range tasksMap {
		tasks = append(tasks, t)
	}
	return tasks
}

func (m *MockTaskRepository) GetByID(userID string, taskID uint64) (models.Task, bool) {
	task, ok := m.tasks[userID][taskID]
	return task, ok
}

func (m *MockTaskRepository) Update(userID string, taskID uint64, updated models.Task) bool {
	if m.tasks[userID] == nil {
		return false
	}
	if _, ok := m.tasks[userID][taskID]; !ok {
		return false
	}
	m.tasks[userID][taskID] = updated
	return true
}

func (m *MockTaskRepository) Delete(userID string, taskID uint64) bool {
	if m.tasks[userID] == nil {
		return false
	}
	if _, ok := m.tasks[userID][taskID]; !ok {
		return false
	}
	delete(m.tasks[userID], taskID)
	return true
}

func setupHandler() *handlers.TaskHandler {
	mockRepo := NewMockTaskRepository()
	service := services.NewTaskService(mockRepo)
	return &handlers.TaskHandler{TaskService: service}
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	handler.TaskService.CreateTask(context.Background(), "user1", dto.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Description",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	handler.GetAllTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp res.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	tasks, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, tasks, 1)
}

func TestTaskHandler_GetTaskByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	created := handler.TaskService.CreateTask(context.Background(), "user1", dto.CreateTaskRequest{
		Title:       "Sample Task",
		Description: "Desc",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", created.ID)}}

	handler.GetTaskByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp res.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(created.ID), data["ID"])
	assert.Equal(t, "Sample Task", data["Title"])
}

func TestTaskHandler_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	newTask := dto.CreateTaskRequest{
		Title:       "New Task",
		Description: "New Description",
	}
	body, _ := json.Marshal(newTask)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp res.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "New Task", data["Title"])
	assert.Equal(t, "New Description", data["Description"])
}

func TestTaskHandler_GetAllTasks_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler.GetAllTasks(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Failed to retrieve user ID")
}

func TestTaskHandler_GetTaskByID_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler.GetTaskByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Failed to retrieve user ID")
}

func TestTaskHandler_GetTaskByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Params = gin.Params{{Key: "id", Value: "not-a-number"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler.GetTaskByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Invalid task ID")
}

func TestTaskHandler_GetTaskByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler.GetTaskByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Task not found")
}

func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte("{invalid-json")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Invalid request body")
}

func TestTaskHandler_CreateTask_ValidationFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	invalidTask := dto.CreateTaskRequest{
		Title:       "",
		Description: "desc",
	}
	body, _ := json.Marshal(invalidTask)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user1")
	c.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Validation failed")
}

func TestTaskHandler_CreateTask_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupHandler()

	newTask := dto.CreateTaskRequest{
		Title:       "New Task",
		Description: "New Description",
	}
	body, _ := json.Marshal(newTask)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTask(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp res.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Failed to retrieve user ID")
}
