package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/hellofresh/health-go/v4"
	healthGrpc "github.com/hellofresh/health-go/v4/checks/grpc"
	healthHttp "github.com/hellofresh/health-go/v4/checks/http"
	healthMongo "github.com/hellofresh/health-go/v4/checks/mongo"
	healthMysql "github.com/hellofresh/health-go/v4/checks/mysql"
	healthRabbitmq "github.com/hellofresh/health-go/v4/checks/rabbitmq"
	"google.golang.org/grpc"
)

func main() {
	// add some checks on instance creation
	healthz, healthzErr := buildHealthWithChecks()
	if healthzErr != nil {
		panic(healthzErr)
	}

	// register some more checks if needed
	regErr := registerAdditionalChecks(healthz)
	if regErr != nil {
		panic(regErr)
	}

	// handlerUsageSample_httpStandard(healthz)

	handlerUsageSample_muxRouter(healthz)

	// handlerFuncUsageSample_chiRouter(healthz)
}

func handlerUsageSample_httpStandard(healthz *health.Health) {
	http.Handle("/healthz", healthz.Handler())
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func handlerUsageSample_muxRouter(healthz *health.Health) {
	router := mux.NewRouter()
	router.Handle("/healthz", healthz.Handler())
	http.ListenAndServe("0.0.0.0:8080", router)
}

func handlerFuncUsageSample_chiRouter(healthz *health.Health) {
	router := chi.NewRouter()
	router.Get("/healthz", healthz.HandlerFunc)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func buildHealthWithChecks() (*health.Health, error) {
	return health.New(
		health.WithChecks(
			health.Config{
				Name:    "http",
				Timeout: 5 * time.Second,
				Check: healthHttp.New(
					healthHttp.Config{
						URL:            ":8080",
						RequestTimeout: 5 * time.Second,
					},
				),
			},

			health.Config{
				Name:    "grpc",
				Timeout: 5 * time.Second,
				Check: healthGrpc.New(
					healthGrpc.Config{
						Target:  ":50001",
						Service: "testService",
						DialOptions: []grpc.DialOption{
							grpc.WithInsecure(),
						},
					},
				),
			},

			health.Config{
				Name:      "mongodb",
				Timeout:   2 * time.Second,
				SkipOnErr: false,
				Check: healthMongo.New(
					healthMongo.Config{
						DSN:               "mongodb://username:password@0.0.0.0:27017/defaultAuthDb?options",
						TimeoutConnect:    5 * time.Second,
						TimeoutDisconnect: 3 * time.Second,
						TimeoutPing:       2 * time.Second,
					},
				),
			},

			health.Config{
				Name:      "mysql",
				Timeout:   2 * time.Second,
				SkipOnErr: false,
				Check: healthMysql.New(
					healthMysql.Config{
						DSN: "username:password@tcp(0.0.0.0:31726)/dbName?charset=utf8",
					},
				),
			},
		),
	)
}

func registerAdditionalChecks(healthz *health.Health) error {
	return healthz.Register(
		health.Config{
			Name:      "rabbitmq",
			Timeout:   2 * time.Second,
			SkipOnErr: false,
			Check: healthRabbitmq.New(
				healthRabbitmq.Config{
					DSN: "amqp://username:password@0.0.0.0:5672/segment?query",
					// Exchange: "", // application health check exchange - default 'health_check'
					// RoutingKey: "", // application health check routing key within health check exchange - default to host name
					// Queue: "", // application health check queue, that binds to the exchange with the routing key - default '<exchange>.<routing-key>'
					ConsumeTimeout: 3 * time.Second,
				},
			),
		},
	)
}
