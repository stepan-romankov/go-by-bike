package controllers

import (
	"github.com/stepan-romankov/go-by-bike/auth"
	"net/http"
)

func (server *HttpAppServer) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := server.AuthStore.GetSession(r)
	delete(session.Values, auth.SESSION_USER_ID)
	session.Save(r, w)
}
