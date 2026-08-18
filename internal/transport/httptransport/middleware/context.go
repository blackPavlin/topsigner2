package middleware

import "context"

type contextKey string

const userAgentContextKey contextKey = "user-agent"

func GetUserAgentFromContext(ctx context.Context) string {
	userAgent, _ := ctx.Value(userAgentContextKey).(string)

	return userAgent
}

func SetUserAgentToContext(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, userAgentContextKey, userAgent)
}
