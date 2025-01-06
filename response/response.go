package response

import (
	"context"
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	Code      int    `json:"code"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

// RespondWithError sends a standardized JSON error response.
func RespondWithError(ctx context.Context, w http.ResponseWriter, statusCode int, message string, err error, traceID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrResponse{
		Code:      statusCode,
		Reason:    err.Error(),
		Message:   message,
		ErrorCode: traceID,
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": response,
	})
}

// RespondWithSuccess sends a standardized JSON success response.
func RespondWithSuccess(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}
