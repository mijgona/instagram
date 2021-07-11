package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
)


var ErrNoAuthentication = errors.New("no authentication")

var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}
type Auth struct {
	ID		int64
	IsAdmin	bool
}

func (c *contextKey) String() string {
	return c.name
}

type IDFunc func (ctx context.Context, token string) (Auth, error)


func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request){
			token := request.Header.Get("Authorization")
			if request.FormValue("token")!=""{token =request.FormValue("token")}

			auth, err := idFunc(request.Context(), token)
			if err != nil {
				log.Print("Authorization() err:", err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(request.Context(), authenticationContextKey, auth)
			request = request.WithContext(ctx)
			handler.ServeHTTP(writer, request)
		})
	}
}
func Authentication(ctx context.Context) (Auth, error){
	if value, ok :=ctx.Value(authenticationContextKey).(Auth); ok {
		return value, nil
	}
	return Auth{
		ID:      0,
		IsAdmin: false,
	}, ErrNoAuthentication
}