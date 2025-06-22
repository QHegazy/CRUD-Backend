package services

import (
	"context"
	"task-backend/internal/dto"
	"task-backend/internal/models"
	my_utils "task-backend/utils"
)

type TaskRepository interface {
	Create(userID string, task models.Task) models.Task
	GetAll(userID string) []models.Task
	GetByID(userID string, taskID uint64) (models.Task, bool)
	Update(userID string, taskID uint64, updated models.Task) bool
	Delete(userID string, taskID uint64) bool
}

type TaskService struct {
	store TaskRepository
}

func NewTaskService(store TaskRepository) *TaskService {
	return &TaskService{store: store}
}

func (s *TaskService) CreateTask(ctx context.Context, userID string, newTask dto.CreateTaskRequest) models.Task {
	task := models.Task{
		Title:       newTask.Title,
		Description: newTask.Description,
		UserID:      userID,
	}

	created := s.store.Create(userID, task)
	created.ID = my_utils.ObfuscateNumbers(created.ID)
	return created
}

func (s *TaskService) GetAllTasks(ctx context.Context, userID string) []models.Task {
	tasks := s.store.GetAll(userID)

	for i := range tasks {
		tasks[i].ID = my_utils.ObfuscateNumbers(tasks[i].ID)
	}
	return tasks
}

func (s *TaskService) GetTaskByID(ctx context.Context, userID string, taskID uint64) (models.Task, bool) {
	realID := my_utils.DeobfuscateNumbers(taskID)
	task, found := s.store.GetByID(userID, realID)
	if found {
		task.ID = my_utils.ObfuscateNumbers(task.ID)
	}
	return task, found
}

func (s *TaskService) UpdateTask(ctx context.Context, userID string, taskID uint64, updateData dto.UpdateTaskRequest) (models.Task, bool) {
	realID := my_utils.DeobfuscateNumbers(taskID)
	existingTask, found := s.store.GetByID(userID, realID)
	if !found {
		return models.Task{}, false
	}

	if updateData.Title != nil {
		existingTask.Title = *updateData.Title
	}
	if updateData.Description != nil {
		existingTask.Description = *updateData.Description
	}

	ok := s.store.Update(userID, realID, existingTask)
	if !ok {
		return models.Task{}, false
	}

	existingTask.ID = my_utils.ObfuscateNumbers(existingTask.ID)
	return existingTask, true
}

func (s *TaskService) DeleteTask(ctx context.Context, userID string, taskID uint64) bool {
	realID := my_utils.DeobfuscateNumbers(taskID)
	return s.store.Delete(userID, realID)
}
