package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/geowa4/user-service/health"
	"github.com/geowa4/user-service/user"
	"github.com/geowa4/user-service/wrappers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	pgx_stdlib "github.com/jackc/pgx/stdlib"
	log "gopkg.in/inconshreveable/log15.v2"
)

// ServiceRouter defines the interface for all modules handling requests.
type ServiceRouter interface {
	HandleRoutes()
}

func initDB() (*sql.DB, error) {
	dbName := "user-service"
	dbLogger := log.New("db", dbName)
	connConfig := pgx.ConnConfig{
		User:     os.Getenv("DB_ENV_POSTGRES_USER"),
		Password: os.Getenv("DB_ENV_POSTGRES_PASSWORD"),
		Host:     os.Getenv("DB_PORT_5432_TCP_ADDR"),
		Database: dbName,
		Logger:   dbLogger,
	}
	config := pgx.ConnPoolConfig{ConnConfig: connConfig}
	pool, err := pgx.NewConnPool(config)
	if err != nil {
		return nil, err
	}

	db, err := pgx_stdlib.OpenFromConnPool(pool)
	if err != nil {
		log.Crit("unable to create connection pool", "error", err)
		os.Exit(1)
	}

	return db, nil
}

func loadRouters(r *mux.Router) {
	db, err := initDB()
	if err != nil {
		log.Crit("could not initialize db pool", "error", err.Error())
		os.Exit(1)
	}
	services := []ServiceRouter{
		user.NewRouter(db, r.PathPrefix("/"+user.GetName()).Subrouter()),
		health.NewRouter(db, r.PathPrefix("/"+health.GetName()).Subrouter()),
	}
	for _, s := range services {
		s.HandleRoutes()
	}
}

func ping(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte("OK"))
}

func main() {
	log.Info(os.Getenv("GITHUB_CLIENT_ID"))
	log.Info(os.Getenv("GITHUB_CLIENT_SECRET"))

	r := mux.NewRouter().StrictSlash(true)
	r.
		Methods("GET").
		Path("/").
		Name("Home").
		Handler(wrappers.Defaults(ping, "Home"))
	loadRouters(r)
	port := "8080"
	log.Info("service starting", "service", "user", "port", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Crit("service failed to start", "service", "user", "error", err)
	}
}
