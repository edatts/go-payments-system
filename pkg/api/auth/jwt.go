package auth

import "context"

const userKey = "userId"

func GetUserIdFromContext(ctx context.Context) (int32, bool) {
	val, ok := ctx.Value(userKey).(int32)
	return val, ok
}
