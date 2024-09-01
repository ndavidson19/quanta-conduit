package utils

import (
    "fmt"
    "net/http"
)

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func NewAppError(code int, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
    }
}

var (
    ErrNotFound   = NewAppError(http.StatusNotFound, "Resource not found")
    ErrBadRequest = NewAppError(http.StatusBadRequest, "Bad request")
    ErrInternal   = NewAppError(http.StatusInternalServerError, "Internal server error")
)