package storage

import (
	"image_redactor/http"
	"sync"
)

type InMemoryTaskStorage struct {
	mu    sync.Mutex
	tasks map[string]*http.Task
}

func NewInMemoryTaskStorage() *InMemoryTaskStorage {
	return &InMemoryTaskStorage{
		tasks: make(map[string]*http.Task),
	}
}

func (s *InMemoryTaskStorage) PostTask(task *http.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *InMemoryTaskStorage) GetTask(id string) (*http.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, exists := s.tasks[id]
	return task, exists
}
