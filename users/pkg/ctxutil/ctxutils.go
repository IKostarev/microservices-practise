package ctxutil

import (
	"context"
)

// GetUserIDFromContext extracts the user ID from the context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("UserID").(int)
	return userID, ok
}

// GetRequestIDFromContext extracts the request ID from the context
func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value("RequestID").(string)
	return requestID, ok
}
