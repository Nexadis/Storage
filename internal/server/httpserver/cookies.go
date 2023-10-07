package httpserver

import (
	"net/http"

	"github.com/Nexadis/Storage/internal/server/features"
	"github.com/Nexadis/Storage/internal/storage"
)

const (
	UserCookie = "Username"
)

func GetUser(r *http.Request) string {
	if features.FeatureEnabled(features.Use_user, r) {
		c, err := r.Cookie(UserCookie)
		if err != nil {
			return storage.DefaultUser
		}
		return c.Value
	}
	return storage.DefaultUser
}
