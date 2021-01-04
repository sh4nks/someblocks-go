package server

import (
	"net/http"
	"someblocks/internal/app"
)

type RequireUser struct {
	app *app.App
}

func (mw *RequireUser) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := mw.app.GetCurrentUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := r.Context()
		ctx = app.WithCurrentUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
