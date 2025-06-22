package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"task-backend/internal/router"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
