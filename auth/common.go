package auth

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
)

const SESSION_USER_ID = "user_id"

type Store struct {
	SessionStore sessions.Store
}

type Error struct {
	Err error
}

func (r *Error) Error() string {
	return r.Err.Error()
}

func (authStore *Store) GetSession(r *http.Request) (*sessions.Session, error) {
	const SessionName = "auth"
	return authStore.SessionStore.Get(r, SessionName)
}

func (authStore *Store) GetUserId(r *http.Request) (uint32, error) {
	session, err := authStore.GetSession(r)
	if err != nil {
		return 0, err
	}
	userId, ok := session.Values[SESSION_USER_ID].(uint32)
	if !ok {
		return 0, &Error{Err: errors.New("Unauthorized")}
	}
	return userId, nil
}
