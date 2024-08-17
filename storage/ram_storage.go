package storage

import (
	"image_redactor/types"
	"sync"
)

type InMemoryTaskStorage struct {
	mu    sync.Mutex
	tasks map[string]*types.Task
}

func NewInMemoryTaskStorage() *InMemoryTaskStorage {
	return &InMemoryTaskStorage{
		tasks: make(map[string]*types.Task),
	}
}

func (s *InMemoryTaskStorage) PostTask(task *types.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *InMemoryTaskStorage) GetTask(id string) (*types.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, exists := s.tasks[id]
	return task, exists
}
