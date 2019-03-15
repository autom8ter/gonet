package gonet

import (
	"github.com/gorilla/sessions"
)

func NewSessionCookieStore(key string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(key))
}