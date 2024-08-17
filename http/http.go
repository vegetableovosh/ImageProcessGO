package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "image_redactor/docs"
	"image_redactor/types"
	"net/http"
	"time"
)

type Server struct {
	storage types.TaskStorage
}

func NewServer(storage types.TaskStorage) *Server {
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
	//taskID := uuid.New().String()
	task := types.NewTask()

	s.storage.PostTask(task)

	// Эмуляция длительной обработки
	go func(taskID string) {
		time.Sleep(5 * time.Second) // Имитация работы
		task, _ := s.storage.GetTask(taskID)
		task.Status = types.StatusReady
		task.Result = fmt.Sprintf("Processed image data for task %s", taskID)
	}(task.ID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID})
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
	if task.Status != types.StatusReady {
		http.Error(w, "Task is not ready yet", http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": task.Result})
}

func CreateAndRunServer(storage types.TaskStorage, addr string) error {
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
