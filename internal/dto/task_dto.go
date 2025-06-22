package dto

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

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

func (r *CreateTaskRequest) Validate() map[string]string {
	err := validate.Struct(r)
	if err == nil {
		return nil
	}
	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		switch field {
		case "Title":
			switch e.Tag() {
			case "required":
				errors["title"] = "Title is required"
			case "min":
				errors["title"] = "Title must be at least 5 characters"
			case "max":
				errors["title"] = "Title must not exceed 100 characters"
			}
		case "Description":
			switch e.Tag() {
			case "required":
				errors["description"] = "Description is required"
			case "min":
				errors["description"] = "Description must be at least 8 characters"
			case "max":
				errors["description"] = "Description must not exceed 250 characters"
			}
		}
	}
	return errors
}
