package services

import (
	"context"
	"task-backend/internal/dto"
	"task-backend/internal/models"
	my_utils "task-backend/utils"
	"testing"
)

type MockTaskStore struct {
	tasks  map[uint64]models.Task
	nextID uint64
}

func NewMockTaskStore() *MockTaskStore {
	return &MockTaskStore{
		tasks:  make(map[uint64]models.Task),
		nextID: 1,
	}
}

func (m *MockTaskStore) Create(userID string, task models.Task) models.Task {
	task.ID = m.nextID
	task.UserID = userID
	m.nextID++
	m.tasks[task.ID] = task
	return task
}

func (m *MockTaskStore) GetAll(userID string) []models.Task {
	var result []models.Task
	for _, t := range m.tasks {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	return result
}

func (m *MockTaskStore) GetByID(userID string, taskID uint64) (models.Task, bool) {
	task, found := m.tasks[taskID]
	if !found || task.UserID != userID {
		return models.Task{}, false
	}
	return task, true
}

func (m *MockTaskStore) Update(userID string, taskID uint64, updated models.Task) bool {
	task, found := m.tasks[taskID]
	if !found || task.UserID != userID {
		return false
	}
	updated.ID = taskID
	updated.UserID = userID
	m.tasks[taskID] = updated
	return true
}

func (m *MockTaskStore) Delete(userID string, taskID uint64) bool {
	task, found := m.tasks[taskID]
	if !found || task.UserID != userID {
		return false
	}
	delete(m.tasks, taskID)
	return true
}

func TestTaskService_CreateTask(t *testing.T) {
	store := NewMockTaskStore()
	service := NewTaskService(store)

	req := dto.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Description here",
	}

	created := service.CreateTask(context.Background(), "user1", req)
	if created.Title != req.Title || created.Description != req.Description || created.UserID != "user1" {
		t.Errorf("CreateTask returned wrong task data: %+v", created)
	}

	if created.ID == 0 {
		t.Error("CreateTask returned task with zero ID")
	}
}

func TestTaskService_GetAllTasks(t *testing.T) {
	store := NewMockTaskStore()
	service := NewTaskService(store)

	store.Create("user1", models.Task{Title: "Task1", Description: "Desc1"})
	store.Create("user2", models.Task{Title: "Task2", Description: "Desc2"})
	store.Create("user1", models.Task{Title: "Task3", Description: "Desc3"})

	tasks := service.GetAllTasks(context.Background(), "user1")

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks for user1, got %d", len(tasks))
	}

	for _, task := range tasks {
		if task.UserID != "user1" {
			t.Errorf("Got task for wrong user: %+v", task)
		}
		if task.ID == 0 {
			t.Error("Task ID should be obfuscated and non-zero")
		}
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	store := NewMockTaskStore()
	service := NewTaskService(store)

	task := store.Create("user1", models.Task{Title: "Task1", Description: "Desc1"})

	obfuscatedID := my_utils.ObfuscateNumbers(task.ID)

	gotTask, found := service.GetTaskByID(context.Background(), "user1", obfuscatedID)
	if !found {
		t.Error("Expected to find task, but did not")
	}

	if gotTask.Title != task.Title || gotTask.Description != task.Description {
		t.Errorf("Got task does not match original. Got %+v, want %+v", gotTask, task)
	}

	if gotTask.ID == task.ID {
		t.Error("Returned task ID is not obfuscated")
	}

	_, found = service.GetTaskByID(context.Background(), "user2", obfuscatedID)
	if found {
		t.Error("Should not find task for wrong user")
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	store := NewMockTaskStore()
	service := NewTaskService(store)

	task := store.Create("user1", models.Task{Title: "Old Title", Description: "Old Desc"})
	obfuscatedID := my_utils.ObfuscateNumbers(task.ID)

	newTitle := "New Title"
	newDesc := "New Desc"
	updateReq := dto.UpdateTaskRequest{
		Title:       &newTitle,
		Description: &newDesc,
	}

	updatedTask, ok := service.UpdateTask(context.Background(), "user1", obfuscatedID, updateReq)
	if !ok {
		t.Fatal("UpdateTask failed")
	}

	if updatedTask.Title != newTitle || updatedTask.Description != newDesc {
		t.Errorf("UpdateTask did not update fields properly: %+v", updatedTask)
	}

	if updatedTask.ID == task.ID {
		t.Error("Updated task ID is not obfuscated")
	}

	_, ok = service.UpdateTask(context.Background(), "user2", obfuscatedID, updateReq)
	if ok {
		t.Error("UpdateTask should fail for wrong user")
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	store := NewMockTaskStore()
	service := NewTaskService(store)

	task := store.Create("user1", models.Task{Title: "Title", Description: "Desc"})
	obfuscatedID := my_utils.ObfuscateNumbers(task.ID)

	ok := service.DeleteTask(context.Background(), "user1", obfuscatedID)
	if !ok {
		t.Error("DeleteTask failed")
	}

	_, found := store.GetByID("user1", task.ID)
	if found {
		t.Error("Task was not deleted")
	}

	ok = service.DeleteTask(context.Background(), "user2", obfuscatedID)
	if ok {
		t.Error("DeleteTask should fail for wrong user")
	}
}
