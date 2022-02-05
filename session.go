package main

import "net/http"

func isLoggedIn(r *http.Request) bool {
	// Проверяем, что клиент прислал куки с id сессии
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// Проверяем, что сессия не истекла и у нас в БД она есть
	_, ok := SessionDB[c.Value]
	return ok
}

func getUser(w http.ResponseWriter, r *http.Request) User {

	user := User{}

	if sid, err:= r.Cookie("session"); err != nil {
		user = User{}
	} else if username, ok := SessionDB[sid.Value]; !ok {
		user = User{}
	} else {
		if user, ok = UserDB[username]; !ok {
			user = User{}
		}
	}

	return user
}