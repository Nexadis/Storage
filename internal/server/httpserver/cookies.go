package httpserver

import (
	"net/http"

	"github.com/Nexadis/Storage/internal/storage"
)

const (
	UserCookie = "Username"
)

func GetUser(r *http.Request) string {
	c, err := r.Cookie(UserCookie)
	if err != nil {
		return storage.DefaultUser
	}
	return c.Value
}
