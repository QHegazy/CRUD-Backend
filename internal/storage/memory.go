package storage

import (
	"sync"

	"task-backend/internal/models"
)

type TaskStore struct {
	mu        sync.RWMutex
	userTasks map[string]map[uint64]models.Task
	counter   uint64
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		userTasks: make(map[string]map[uint64]models.Task),
		counter:   0,
	}
}

func (s *TaskStore) GetAll(userID string) []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasksMap, exists := s.userTasks[userID]
	if !exists {
		return []models.Task{}
	}

	tasks := make([]models.Task, 0, len(tasksMap))
	for _, task := range tasksMap {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *TaskStore) GetByID(userID string, taskID uint64) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasksMap, exists := s.userTasks[userID]
	if !exists {
		return models.Task{}, false
	}
	task, ok := tasksMap[taskID]
	return task, ok
}

func (s *TaskStore) Create(userID string, task models.Task) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	task.ID = s.counter
	task.UserID = userID

	if _, exists := s.userTasks[userID]; !exists {
		s.userTasks[userID] = make(map[uint64]models.Task)
	}
	s.userTasks[userID][task.ID] = task
	return task
}

func (s *TaskStore) Update(userID string, taskID uint64, updated models.Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasksMap, exists := s.userTasks[userID]
	if !exists {
		return false
	}

	if _, exists := tasksMap[taskID]; !exists {
		return false
	}

	updated.ID = taskID
	updated.UserID = userID
	tasksMap[taskID] = updated
	return true
}

func (s *TaskStore) Delete(userID string, taskID uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasksMap, exists := s.userTasks[userID]
	if !exists {
		return false
	}

	if _, exists := tasksMap[taskID]; !exists {
		return false
	}
	delete(tasksMap, taskID)
	return true
}
