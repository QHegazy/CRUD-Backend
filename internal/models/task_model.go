package models

type Task struct {
	ID          uint64
	UserID      string
	Title       string
	Description string
}
