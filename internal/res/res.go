package res

type ErrorResponse struct {
	Message string `json:"message"`
	Error   any    `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
