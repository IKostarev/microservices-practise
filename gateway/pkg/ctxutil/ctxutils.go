package ctxutil

import (
	"context"
)

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("UserID").(int)
	return userID, ok
}

func SetUserIDToContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, "UserID", userID)
}

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value("RequestID").(string)
	return requestID, ok
}

func SetRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "RequestID", requestID)
}
