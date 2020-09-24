package controllers

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/stepan-romankov/go-by-bike/auth"
	"github.com/stepan-romankov/go-by-bike/middlewares"
	"log"
	"net"
	"net/http"
)

type HttpAppServer struct {
	DB        *pg.DB
	Router    *mux.Router
	AuthStore auth.Store
}

func (server *HttpAppServer) initializeRoutes(authStore auth.Store) *mux.Router {
	routes := mux.NewRouter()
	authMdw := middlewares.NewAuthMiddleware(authStore)

	// TODO: auth middleware builder (inject store)
	routes.HandleFunc("/signup", server.Signup).Methods(http.MethodPost)
	routes.HandleFunc("/logon", server.Logon).Methods(http.MethodPost)
	routes.HandleFunc("/logout", authMdw.Wrap(server.Logout)).Methods(http.MethodPost)
	routes.HandleFunc("/bikes", server.Bikes).Methods(http.MethodGet)
	routes.HandleFunc("/rent", authMdw.Wrap(server.Rent)).Methods(http.MethodPost)
	routes.HandleFunc("/return", authMdw.Wrap(server.Return)).Methods(http.MethodPost)
	return routes
}

func (server *HttpAppServer) Initialize(db *pg.DB, authStore auth.Store) {
	server.DB = db

	server.AuthStore = authStore
	server.Router = server.initializeRoutes(authStore)

}

func (server *HttpAppServer) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("Listening to port %d\n", listener.Addr().(*net.TCPAddr).Port)
	log.Fatal(http.Serve(listener, server.Router))
}
