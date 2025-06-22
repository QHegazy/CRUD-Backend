package storage

import (
	"sync"

	"task-backend/internal/models"
)

type TaskStore struct {
	tasks map[string]models.Task
	mu    sync.RWMutex
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]models.Task),
	}
}

func (s *TaskStore) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *TaskStore) GetByID(id string) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *TaskStore) Create(task models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *TaskStore) Update(id string, task models.Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return false
	}
	s.tasks[id] = task
	return true
}

func (s *TaskStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return false
	}
	delete(s.tasks, id)
	return true
}
