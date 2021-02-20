package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/hellofresh/health-go/v4"
	healthMysql "github.com/hellofresh/health-go/v4/checks/mysql"
)

func main() {
	// add some checks on instance creation
	healtz, healtzErr := buildHealthWithChecks()
	if healtzErr != nil {
		panic(healtzErr)
	}

	// register some more checks if needed
	regErr := registerAdditionalChecks(healtz)
	if regErr != nil {
		panic(regErr)
	}

	// handlerUsageSample_httpStandard(healtz)

	handlerUsageSample_muxRouter(healtz)

	// handlerFuncUsageSample_chiRouter(healtz)
}

func handlerUsageSample_httpStandard(healtz *health.Health) {
	http.Handle("/healtz", healtz.Handler())
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func handlerUsageSample_muxRouter(healtz *health.Health) {
	router := mux.NewRouter()
	router.Handle("/healtz", healtz.Handler())
	http.ListenAndServe("0.0.0.0:8080", router)
}

func handlerFuncUsageSample_chiRouter(healtz *health.Health) {
	router := chi.NewRouter()
	router.Get("/healtz", healtz.HandlerFunc)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func buildHealthWithChecks() (*health.Health, error) {
	return health.New(
		health.WithChecks(
			health.Config{
				Name:      "rabbitmq",
				Timeout:   5 * time.Second,
				SkipOnErr: true,
				Check: func(ctx context.Context) error {
					// rabbitmq health check implementation goes here
					return nil
				},
			},

			health.Config{
				Name: "mongodb",
				Check: func(ctx context.Context) error {
					// mongo_db health check implementation goes here
					return nil
				},
			},
		),
	)
}

func registerAdditionalChecks(healtz *health.Health) error {
	return healtz.Register(health.Config{
		Name:      "mysql",
		Timeout:   2 * time.Second,
		SkipOnErr: false,
		Check: healthMysql.New(
			healthMysql.Config{
				DSN: "test:test@tcp(0.0.0.0:31726)/test?charset=utf8",
			},
		),
	})
}
