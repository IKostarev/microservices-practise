package ctxutil

import (
	"context"
	"google.golang.org/grpc/metadata"
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

func SetRequestIdFromContextToGrpc(ctx context.Context) context.Context {
	requestID, ok := GetRequestIDFromContext(ctx)
	if !ok {
		return ctx
	}

	md := metadata.Pairs("requestId", requestID)
	return metadata.NewOutgoingContext(ctx, md)
}

func SetRequestIdFromGrpcToContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	requestIds, ok := md["requestid"]
	if ok && len(requestIds) > 0 {
		ctx = SetRequestIDToContext(ctx, requestIds[0])
	}

	return ctx
}
