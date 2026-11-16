package models

// модель ошибки для ответа апи
type ErrorResponse struct {
	Error string `json:"error"`
}

// модель успешного ответа апи
type SuccessResponse struct {
	Message string `json:"message"`
}
