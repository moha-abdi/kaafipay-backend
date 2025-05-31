package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

// RespondWithJSON writes a JSON response with the given status code and payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error encoding response", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError writes a JSON error response with the given status code and error details
func RespondWithError(w http.ResponseWriter, code int, errorCode string, message string, details interface{}) {
	errResp := ErrorResponse{}
	errResp.Error.Code = errorCode
	errResp.Error.Message = message
	errResp.Error.Details = details

	response, _ := json.Marshal(errResp)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.WriteHeader(code)
	w.Write(response)
}
