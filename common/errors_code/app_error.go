package errors_code

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ITEM_EXISTS       = &AppError{1001, "Item exists", http.StatusBadRequest}
	ITEM_NO_EXISTS    = &AppError{1002, "Item no exists", http.StatusBadRequest}
	TYPE_ITEM_INVALID = &AppError{1002, "Type item invalid (MAIN, SIDE, DRINK)", http.StatusBadRequest}

	INTERNAL  = &AppError{9999, "Internal server error", http.StatusInternalServerError}
	NOT_FOUND = &AppError{9001, "Resource not found", http.StatusNotFound}
)
