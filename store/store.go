package store

import (
	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitStore(secret string) {
	Store = sessions.NewCookieStore([]byte(secret))
}
