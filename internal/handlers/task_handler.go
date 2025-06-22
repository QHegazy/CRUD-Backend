package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all tasks"})
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get task by ID", "id": id})
}

func CreateTask(c *gin.Context) {

	c.JSON(http.StatusCreated, gin.H{"message": "Task created", "task": "task"})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Task updated", "id": id})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted", "id": id})
}
