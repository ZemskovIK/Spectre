package lib

import "context"

type CtxKey string

var UserAccessLevelKey CtxKey = "userAccessLevel"
var UserIDKey CtxKey = "userID"

func FetchAccessLevelFromCtx(ctx context.Context) (int, bool) {
	alRaw := ctx.Value(UserAccessLevelKey)
	al, ok := alRaw.(float64)
	if !ok {
		return -1, false
	}
	return int(al), true
}
