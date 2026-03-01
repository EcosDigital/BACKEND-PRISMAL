package middlewares

import "context"

type AuthUser struct {
	ID     int
	Email  string
	Nombre string
	RolID  int
	Rol    string
}

type ctxKey string

const userCtxKey ctxKey = "user"

func UserFromCtx(ctx context.Context) (AuthUser, bool) {
	user, ok := ctx.Value(userCtxKey).(AuthUser)
	return user, ok
}
