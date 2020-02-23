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
	fmt.Println("Opening DB connection")
	// For brevity, error check is being omitted here.
	db, _ := sql.Open("mysql", "user:password@/dbname")
	defer db.Close()

	fmt.Println("Preparing router")
	router := mux.NewRouter()
	router.Handle("/healthcheck", healthcheck.Handler(

		// WithTimeout allows you to set a max overall timeout.
		healthcheck.WithTimeout(5*time.Second),

		// Checkers fail the status in case of any error.
		healthcheck.WithChecker(
			"heartbeat", checkers.Heartbeat("$PROJECT_PATH/heartbeat"),
		),

		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
			),
		),

		// Observers do not fail the status in case of error.
		healthcheck.WithObserver(
			"diskspace", checkers.DiskSpace("/var/log", 90),
		),
	))

	fmt.Println("Starting REST server")
	http.ListenAndServe(":8080", router)
}
