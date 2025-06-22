package storage

import (
	"testing"

	"task-backend/internal/models"
)

func TestTaskStore_CreateAndGetByID(t *testing.T) {
	store := NewTaskStore()

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	created := store.Create("user1", task)

	if created.ID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
	if created.UserID != "user1" {
		t.Errorf("Expected UserID to be 'user1', got %s", created.UserID)
	}

	retrieved, found := store.GetByID("user1", created.ID)
	if !found {
		t.Error("Expected to find created task")
	}
	if retrieved.Title != task.Title || retrieved.Description != task.Description {
		t.Errorf("Retrieved task doesn't match created task: %+v", retrieved)
	}
}

func TestTaskStore_GetAll(t *testing.T) {
	store := NewTaskStore()

	store.Create("user1", models.Task{Title: "Task 1"})
	store.Create("user1", models.Task{Title: "Task 2"})
	store.Create("user2", models.Task{Title: "Task X"})

	tasksUser1 := store.GetAll("user1")
	if len(tasksUser1) != 2 {
		t.Errorf("Expected 2 tasks for user1, got %d", len(tasksUser1))
	}

	tasksUser2 := store.GetAll("user2")
	if len(tasksUser2) != 1 {
		t.Errorf("Expected 1 task for user2, got %d", len(tasksUser2))
	}

	tasksUser3 := store.GetAll("user3")
	if len(tasksUser3) != 0 {
		t.Errorf("Expected 0 tasks for unknown user, got %d", len(tasksUser3))
	}
}

func TestTaskStore_Update(t *testing.T) {
	store := NewTaskStore()

	task := store.Create("user1", models.Task{Title: "Old Title", Description: "Old Desc"})

	updated := models.Task{Title: "New Title", Description: "New Desc"}
	success := store.Update("user1", task.ID, updated)
	if !success {
		t.Error("Expected update to succeed")
	}

	afterUpdate, found := store.GetByID("user1", task.ID)
	if !found {
		t.Fatal("Updated task not found")
	}
	if afterUpdate.Title != "New Title" || afterUpdate.Description != "New Desc" {
		t.Errorf("Task not updated properly: %+v", afterUpdate)
	}

	failed := store.Update("user1", 9999, updated)
	if failed {
		t.Error("Expected update to fail for non-existent task ID")
	}

	failed = store.Update("user2", task.ID, updated)
	if failed {
		t.Error("Expected update to fail for wrong user")
	}
}

func TestTaskStore_Delete(t *testing.T) {
	store := NewTaskStore()

	task := store.Create("user1", models.Task{Title: "Task to delete"})

	deleted := store.Delete("user1", task.ID)
	if !deleted {
		t.Error("Expected delete to succeed")
	}

	_, found := store.GetByID("user1", task.ID)
	if found {
		t.Error("Expected task to be deleted")
	}


	deletedAgain := store.Delete("user1", task.ID)
	if deletedAgain {
		t.Error("Expected second delete to fail")
	}

	task2 := store.Create("user1", models.Task{Title: "Another task"})
	deletedWrongUser := store.Delete("user2", task2.ID)
	if deletedWrongUser {
		t.Error("Expected delete to fail for wrong user")
	}
}
