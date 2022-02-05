package main

import (
	"net/http"
	"time"
)

// Устанавливаем максимальную длительность сессии
var sessionLengthInHours = time.Hour * 24

func isLoggedIn(r *http.Request) bool {
	// Проверяем, что клиент прислал куки с id сессии
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// Проверяем, что сессия не истекла и у нас в БД она есть
	s, ok := sessionDB[c.Value]
	if !ok {
		return false
	} else if s.CreationTime.Add(sessionLengthInHours).Before(time.Now()) {
		delete(sessionDB, c.Value)
		return false
	} else {
		return true
	}
}

func getUser(w http.ResponseWriter, r *http.Request) user {

	u := user{}

	// Проверяем, что есть куки с сессией
	c, err:= r.Cookie("session")
	if err != nil {
		return user{}
	}

	// Проверяем, что такая сессия есть в нашей БД
	s, ok := sessionDB[c.Value]
	if !ok {
		return user{}
	}

	// Проверяем, что сессия не устарела
	if s.CreationTime.Add(sessionLengthInHours).Before(time.Now()) {
		delete(sessionDB, c.Value)
		return user{}
	}

	// Проверяем, что есть соответствующий пользователь в БД
	if u, ok = userDB[s.Username]; !ok {
		return user{}
	}
	return u
}