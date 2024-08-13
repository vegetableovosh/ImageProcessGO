package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "image_redactor/docs"
	"net/http"
	"time"
)

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
)

type Task struct {
	ID     string     `json:"id"`
	Status TaskStatus `json:"status"`
	Result string     `json:"result"`
}

type TaskStorage interface {
	PostTask(task *Task)
	GetTask(id string) (*Task, bool)
}

type Server struct {
	storage TaskStorage
}

func NewServer(storage TaskStorage) *Server {
	return &Server{storage: storage}
}

func (s *Server) taskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	task, exists := s.storage.GetTask(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary Create a new task
// @Description Create a new task and return the task ID
// @Success 201 {object} map[string]string
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [post]
func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := uuid.New().String()
	task := &Task{
		ID:     taskID,
		Status: StatusInProgress,
	}

	s.storage.PostTask(task)

	// Эмуляция длительной обработки
	go func(taskID string) {
		time.Sleep(5 * time.Second) // Имитация работы
		task, _ := s.storage.GetTask(taskID)
		task.Status = StatusReady
		task.Result = fmt.Sprintf("Processed image data for task %s", taskID)
	}(taskID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

// @Summary Get task status
// @Description Get the current status of a task by its ID
// @Param taskID path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {string} string "Task Not Found"
// @Router /status/{taskID} [get]
func (s *Server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	task, exists := s.storage.GetTask(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": string(task.Status)})
}

// @Summary Get task result
// @Description Get the result of a task by its ID if it's ready
// @Param taskID path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {string} string "Task Not Found"
// @Failure 202 {string} string "Task Not Ready"
// @Router /result/{taskID} [get]
func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	task, exists := s.storage.GetTask(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	if task.Status != StatusReady {
		http.Error(w, "Task is not ready yet", http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": task.Result})
}

func CreateAndRunServer(storage TaskStorage, addr string) error {
	server := NewServer(storage)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/task", server.createTaskHandler)
	r.Get("/status/{taskID}", server.getStatusHandler)
	r.Get("/result/{taskID}", server.getResultHandler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
