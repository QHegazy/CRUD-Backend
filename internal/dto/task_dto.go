package dto

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=5,max=100"`
	Description string `json:"description" validate:"required,min=8,max=250"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=5,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=8,max=250"`
}

type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
