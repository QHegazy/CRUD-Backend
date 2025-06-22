package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"task-backend/internal/handlers"
	"task-backend/internal/router"
	"task-backend/internal/services"
	"task-backend/internal/storage"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	store := storage.NewTaskStore()

	taskService := services.NewTaskService(store)

	taskHandler := &handlers.TaskHandler{TaskService: taskService}

	handler := router.RegisterRoutes(taskHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
