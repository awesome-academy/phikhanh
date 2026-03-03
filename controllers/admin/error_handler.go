package admin

import (
	"errors"
	"phikhanh/utils"
)

// formatErrorMessage - Extract safe error message from ServiceError or return generic message
func formatErrorMessage(err error) string {
	var svcErr *utils.ServiceError
	if errors.As(err, &svcErr) {
		// Only show Message for client errors (4xx), hide DB details for 5xx
		if svcErr.StatusCode >= 400 && svcErr.StatusCode < 500 {
			return svcErr.Message
		}
		// For 5xx errors, return generic message
		return "An error occurred while processing your request"
	}
	return "An unexpected error occurred"
}
