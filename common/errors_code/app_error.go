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

	ORDER_EXISTS         = &AppError{1003, "Order exists", http.StatusBadRequest}
	ORDER_NO_EXISTS      = &AppError{1004, "Order no exists", http.StatusBadRequest}
	STATUS_ORDER_INVALID = &AppError{1005, "Status Order invalid (PENDING, CANCELLED, CONFIRMED, COMPLETED)", http.StatusBadRequest}
	ORDER_NO_PENDING     = &AppError{1006, "Order no pending", http.StatusBadRequest}
	ORDER_NO_CONFIRM     = &AppError{1007, "Order no confirm", http.StatusBadRequest}

	USER_EXISTS       = &AppError{1006, "User exists", http.StatusBadRequest}
	USER_NO_EXISTS    = &AppError{1007, "User no exists", http.StatusBadRequest}
	ROLE_USER_INVALID = &AppError{1008, "Role user invalid (ADMIN, CUSTOMER)", http.StatusBadRequest}

	INTERNAL  = &AppError{9999, "Internal server error", http.StatusInternalServerError}
	NOT_FOUND = &AppError{9001, "Resource not found", http.StatusNotFound}
)
