package router

import (
	"io"
	"log"
	"net/http"
	"os"
	"task-backend/internal/handlers"
	"task-backend/internal/middlewares"
	"task-backend/internal/res"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(taskHandler *handlers.TaskHandler) http.Handler {
	r := gin.New()

	f, err := os.OpenFile("gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open gin log file: %v", err)
	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsOrigin := os.Getenv("CORS")
	if corsOrigin == "" {
		log.Println("CORS env not found, defaulting to http://localhost:3000")
		corsOrigin = "http://localhost:3000"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: corsOrigin != "*",
	}))

	taskGroup := r.Group("/tasks", middlewares.AuthMiddleware())
	{
		taskGroup.GET("", taskHandler.GetAllTasks)
		taskGroup.GET("/:id", taskHandler.GetTaskByID)
		taskGroup.POST("", taskHandler.CreateTask)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)
		taskGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, res.ErrorResponse{
			Message: "Route not found",
			Status:  http.StatusNotFound,
			Error:   "invalid route",
		})
	})

	return r
}
