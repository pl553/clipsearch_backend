package dtos

// swagger:model JsendEmptySuccessResponse
type JsendEmptySuccessResponse struct {
	// Set to "success"
	Status string `json:"status" example:"success"`
	Data   any    `json:"data" type:"null"`
}

// When an API call failed due to invalid parameters
// swagger:model JsendFailResponse
type JsendFailResponse struct {
	// Set to "fail"
	Status string            `json:"status" example:"fail"`
	Data   map[string]string `json:"data" example:"bar:invalid,foo:also invalid"`
}

// When an API call failed due to an error on the server
// swagger:model JsendErrorResponse
type JsendErrorResponse struct {
	// Set to "error"
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"An internal error has occurred"`
}

func NewJsendEmptySuccessResponse() JsendEmptySuccessResponse {
	return JsendEmptySuccessResponse{
		Status: "success",
		Data:   nil,
	}
}

func NewJsendFailResponse(data map[string]string) JsendFailResponse {
	return JsendFailResponse{
		Status: "fail",
		Data:   data,
	}
}

func NewJsendErrorResponse(message string) JsendErrorResponse {
	return JsendErrorResponse{
		Status:  "error",
		Message: message,
	}
}
