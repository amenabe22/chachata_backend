package middlewares

import (
	"context"
	"net/http"

	"github.com/amenabe22/chachata_backend/graph/jwt"
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/amenabe22/chachata_backend/graph/setup"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

var db = setup.SetupModels()

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := header
			usrID, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// create user and check if user exists in db
			user := &model.User{ID: usrID}
			db.Find(&user)
			// println(user.ID, "HERE")
			// println(uint(finID), "TST")
			// id, err := model.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			// put it in context
			user.ID = usrID
			ctx := context.WithValue(r.Context(), userCtxKey, &user)
			// println(user.Name)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.

func ForContext(ctx context.Context) *model.User {
	// println(ctx.Err().Error())
	// TODO: FIX server is crashing when request is sent with no header
	raw, _ := ctx.Value(userCtxKey).(**model.User)
	println(*raw, "here")
	return *raw
}
