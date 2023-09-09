package server

import "net/http"

const (
	UserCookie  = "Username"
	defaultUser = "default"
)

func GetUser(r *http.Request) string {
	c, err := r.Cookie(UserCookie)
	if err != nil {
		return defaultUser
	}
	return c.String()
}
