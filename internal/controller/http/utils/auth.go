package utils

import "context"

var (
	usernameContextKey = struct{}{}
)

func AddUsernameToContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameContextKey, username)
}

func GetUsernameFromContext(ctx context.Context) string {
	v := ctx.Value(usernameContextKey)
	if v == nil {
		return ""
	}

	if username, ok := v.(string); ok {
		return username
	}

	return ""
}
