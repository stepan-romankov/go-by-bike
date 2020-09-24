package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pgext"
	"github.com/gorilla/sessions"
	"github.com/stepan-romankov/go-by-bike/auth"
	"github.com/stepan-romankov/go-by-bike/controllers"
	"github.com/stepan-romankov/go-by-bike/db"
	"log"
	"math/rand"
	"os"
	"strconv"
)

var server = controllers.HttpAppServer{}

func main() {
	postgresUrl, exists := os.LookupEnv("POSTGRES_URL")
	if !exists {
		log.Panicf("POSTGRES_URL env var is not defined")
	}

	postgresOpts, err := pg.ParseURL(os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	serverAddress, exists := os.LookupEnv("APP_ADDRESS")
	if !exists {
		serverAddress = ":8080"
	}

	migrationsPth, exists := os.LookupEnv("POSTGRES_MIGRATIONS_PATH")
	if !exists {
		migrationsPth = "file://db/migrations"
	}

	con := pg.Connect(postgresOpts)
	con.AddQueryHook(pgext.DebugHook{Verbose: true})

	_ = db.Migrate(postgresUrl, migrationsPth)

	fixtures, _ := strconv.ParseBool(os.Getenv("FIXTURES"))
	if fixtures {
		db.Fixtures(con)
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if len(sessionKey) == 0 {
		sessionKey = fmt.Sprintf("%d", rand.Int())
	}

	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	authStore := auth.Store{SessionStore: sessionStore}
	server.Initialize(con, authStore)
	server.Run(serverAddress)
}
