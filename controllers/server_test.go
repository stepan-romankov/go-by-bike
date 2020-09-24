package controllers

import (
	"bytes"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pgext"
	"github.com/gorilla/sessions"
	"github.com/stepan-romankov/go-by-bike/auth"
	"github.com/stepan-romankov/go-by-bike/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var Server HttpAppServer

func TestMain(m *testing.M) {
	postgresUrl := os.Getenv("POSTGRES_TEST_URL")
	postgresOpts, err := pg.ParseURL(postgresUrl)
	if err != nil {
		panic(err)
	}

	con := pg.Connect(postgresOpts)
	con.AddQueryHook(pgext.DebugHook{Verbose: true})
	db.Cleanup(con)

	migrationsPth, exists := os.LookupEnv("POSTGRES_MIGRATIONS_PATH")
	if !exists {
		migrationsPth = "file://db/migrations"
	}
	db.Migrate(postgresUrl, migrationsPth)

	db.Fixtures(con)

	sessionStore := sessions.NewCookieStore([]byte("s"))
	authStore := auth.Store{SessionStore: sessionStore}
	Server.Initialize(con, authStore)
	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	Server.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func login(login string, password string) *http.Cookie {
	body := bytes.NewBuffer([]byte(`{"login":"` + login + `", "password": "` + password + `"}`))
	req, _ := http.NewRequest(http.MethodPost, "/logon", body)
	response := executeRequest(req)
	cookieStr := response.Header()["Set-Cookie"]

	request := &http.Request{Header: http.Header{"Cookie": cookieStr}}

	// Extract the dropped cookie from the request.
	cookie, _ := request.Cookie("auth")
	return cookie
}
