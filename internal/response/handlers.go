package response

// ErrorResponse представляет стандартный формат ответа при ошибке
// @Description Стандартный ответ при ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}
