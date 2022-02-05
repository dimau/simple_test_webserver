package main

import "net/http"

func isLoggedIn(r *http.Request) bool {
	// Проверяем, что клиент прислал куки с id сессии
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// Проверяем, что сессия не истекла и у нас в БД она есть
	_, ok := sessionDB[c.Value]
	return ok
}

func getUser(w http.ResponseWriter, r *http.Request) user {

	u := user{}

	if sid, err:= r.Cookie("session"); err != nil {
		u = user{}
	} else if username, ok := sessionDB[sid.Value]; !ok {
		u = user{}
	} else {
		if u, ok = userDB[username]; !ok {
			u = user{}
		}
	}

	return u
}