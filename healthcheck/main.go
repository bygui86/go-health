package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
	"github.com/etherlabsio/healthcheck/checkers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Open DB connection")
	db, dbErr := sql.Open("mysql", "user:password@/dbname")
	if dbErr != nil {
		// WARN: for sake of demonstrating healthcheck, error is just logged!
		fmt.Println("Starting REST server")
		// panic(dbErr)
	}

	fmt.Println("Prepare HTTP router")
	router := mux.NewRouter()
	router.Handle("/healtz", getHealtzHandler(db))

	fmt.Println("Start REST server")
	http.ListenAndServe("0.0.0.0:8080", router)

	if db != nil {
		fmt.Println("Close DB connection")
		db.Close()
	}
}

func getHealtzHandler(db *sql.DB) http.Handler {
	return healthcheck.Handler(
		// WithTimeout allows you to set a max overall timeout
		healthcheck.WithTimeout(5*time.Second),

		// Checkers fail the status in case of any error
		healthcheck.WithChecker(
			"heartbeat", checkers.Heartbeat("./heartbeat"),
		),
		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
			),
		),

		// Observers do not fail the status in case of error
		healthcheck.WithObserver(
			"diskspace", checkers.DiskSpace("/var/log", 90),
		),
		healthcheck.WithObserver(
			"customSvc", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return fmt.Errorf("custom service failure")
				},
			),
		),
	)
}
