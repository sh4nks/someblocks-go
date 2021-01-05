package middleware

import (
	"net/http"
	"someblocks/internal/app"
	"someblocks/internal/context"
	"someblocks/internal/models"
	"strings"
)

func NewUserMiddleware(app *app.App, userService *models.UserService) *User {
	return &User{
		app:         app,
		userService: userService,
	}
}

type User struct {
	app         *app.App
	userService *models.UserService
}

func (mw *User) CurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// If the user is requesting a static asset we will not need to lookup
		// the current user so we skip doing that.
		if strings.HasPrefix(path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		userId := mw.app.Session.GetInt(r.Context(), "userId")
		if userId == 0 {
			next.ServeHTTP(w, r)
			return
		}

		user := mw.userService.GetById(userId)
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithCurrentUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// LoginRequired assumes that CurrentUser middleware has already been run
// otherwise it will no work correctly.
func (mw *User) LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.CurrentUser(r.Context())
		if user == nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
