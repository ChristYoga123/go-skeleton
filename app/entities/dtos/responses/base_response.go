package responses

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// NewSuccessResponse membuat success response
// Gunakan http.StatusOK (200), http.StatusCreated (201), dll untuk parameter code
// Contoh: NewSuccessResponse("Success", data, http.StatusOK)
func NewSuccessResponse(message string, data any, code int) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse membuat error response
// Gunakan http.StatusInternalServerError (500), http.StatusBadRequest (400), dll untuk parameter code
// Contoh: NewErrorResponse("Error", nil, http.StatusBadRequest)
func NewErrorResponse(message string, data any, code int) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
