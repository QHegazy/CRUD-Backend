package res

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   any    `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}
