package handlers

import (
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitStore(s *sessions.CookieStore) {
	store = s
}

// 添加一个获取 store 的函数
func GetStore() *sessions.CookieStore {
	return store
}
