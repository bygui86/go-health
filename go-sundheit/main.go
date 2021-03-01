package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/AppsFlyer/go-sundheit"
	"github.com/pkg/errors"

	"github.com/AppsFlyer/go-sundheit/checks"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
)

func main() {
	// *** create a new health instance

	// simple
	// health := gosundheit.New()

	// with custom listeners
	health := gosundheit.New(
		gosundheit.WithCheckListeners(&checkEventsLogger{}),
		gosundheit.WithHealthListeners(&checkHealthLogger{}),
	)

	// WARN: don't ignore errors!!!

	// *** define checks

	// HTTP built-in check

	// create the HTTP check for the dependency
	// fail fast if you mis-configured the URL.
	httpCheck, httpErr := checks.NewHTTPCheck(
		checks.HTTPCheckConfig{
			CheckName: "httpbin.url.check",
			Timeout:   1 * time.Second,
			// dependency you're checking - use your own URL here...
			// URL: "http://httpbin.org/status/200,300", // this URL will fail 50% of the times
			URL: "http://localhost:8080/api",
		},
	)
	if httpErr != nil {
		fmt.Println("[ERROR] Failed to create http check: ", httpErr)
		os.Exit(501) // your call...
	}

	// alternatively panic when creating a check fails
	// httpCheck := checks.Must(checks.NewHTTPCheck(httpCheckConf))

	httpRegErr := health.RegisterCheck(
		&gosundheit.Config{
			Check:           httpCheck,
			InitialDelay:    time.Second,      // run 1st time after 1 sec
			ExecutionPeriod: 10 * time.Second, // run every 10 sec
		},
	)
	if httpRegErr != nil {
		fmt.Println("[ERROR] Failed to register check: ", httpRegErr.Error())
		os.Exit(501) // or whatever
	}

	// DNS built-in check
	// Schedule a host resolution check for `example.com`, requiring at least one results
	dnsRegErr := health.RegisterCheck(
		&gosundheit.Config{
			Check: checks.NewHostResolveCheck(
				"example.com",
				200*time.Millisecond,
				1,
			),
			ExecutionPeriod: 10 * time.Second,
		},
	)
	if dnsRegErr != nil {
		fmt.Println("[ERROR] Failed to register check: ", dnsRegErr.Error())
		os.Exit(501) // or whatever
	}

	// Ping built-in check(s)
	// use it as a DB ping check (sql.DB implements the Pinger interface)
	// WARN: commented as "db" is nil and gosunheit generates a nil pointer dereference
	// db, _ := sql.Open("mysql", "user:password@/dbname") // WARN: for sake of demonstrating the health check, error is ignored
	// dbCheck, dbErr := checks.NewPingCheck("db.check", db, 100*time.Millisecond)
	// if dbErr != nil {
	// 	fmt.Println("[ERROR] Failed to create db check: ", dbErr)
	// 	os.Exit(501) // your call...
	// }
	// dbRegErr := health.RegisterCheck(
	// 	&gosundheit.Config{
	// 		Check:           dbCheck,
	// 		ExecutionPeriod: 10 * time.Second,
	// 	},
	// )
	// if dbRegErr != nil {
	// 	fmt.Println("[ERROR] Failed to register check: ", dbRegErr)
	// 	os.Exit(501) // or whatever
	// }

	// use the ping check to test a generic connection
	pingCheck, pingErr := checks.NewPingCheck(
		// "example.com.reachable",
		"localhost.reachable",
		// checks.NewDialPinger("tcp", "example.com"),
		checks.NewDialPinger("tcp", "localhost:8080"),
		time.Second,
	)
	if pingErr != nil {
		fmt.Println("[ERROR] Failed to create ping check: ", pingErr)
		os.Exit(501) // your call...
	}
	pingerErr := health.RegisterCheck(
		&gosundheit.Config{
			Check:           pingCheck,
			ExecutionPeriod: 10 * time.Second,
		},
	)
	if pingerErr != nil {
		fmt.Println("[ERROR] Failed to register check: ", pingerErr)
		os.Exit(501) // or whatever
	}

	// Custom built-in check(s)
	// use the CustomCheck struct
	customStructErr := health.RegisterCheck(
		&gosundheit.Config{
			Check: &checks.CustomCheck{
				CheckName: "lottery.check.struct",
				CheckFunc: lotteryCheck,
			},
			InitialDelay:    3 * time.Second,
			ExecutionPeriod: 5 * time.Second,
		},
	)
	if customStructErr != nil {
		fmt.Println("[ERROR] Failed to register check: ", customStructErr)
		os.Exit(501) // or whatever
	}

	// implement the Check interface
	checkIntErr := health.RegisterCheck(
		&gosundheit.Config{
			Check: &Lottery{
				name:        "lottery.check.interface",
				probability: 0.3,
			},
			InitialDelay:    1 * time.Second,
			ExecutionPeriod: 30 * time.Second,
		},
	)
	if checkIntErr != nil {
		fmt.Println("[ERROR] Failed to register check: ", checkIntErr)
		os.Exit(501) // or whatever
	}

	// *** define more checks...

	// *** register endpoints
	// api
	http.Handle("/api", apiHandler())
	// health endpoint
	// http.Handle("/admin/health.json", healthhttp.HandleHealthJSON(health))
	http.Handle("/healthz", healthhttp.HandleHealthJSON(health))

	// *** serve
	fmt.Println(fmt.Sprintf(
		"[ERROR] HTTP server failed: %s",
		http.ListenAndServe("0.0.0.0:8080", nil).Error(),
	))
}

// ***

// custom checks

// custom struct

func lotteryCheck() (details interface{}, err error) {
	lottery := rand.Float32()
	details = fmt.Sprintf("lottery=%f", lottery)
	if lottery < 0.5 {
		err = errors.New("[ERROR] Sorry, I failed")
	}
	return
}

// checks.Check interface

type Lottery struct {
	name        string
	probability float32
}

func (l *Lottery) Execute() (details interface{}, err error) {
	return lotteryCheck()
}

func (l *Lottery) Name() string {
	return l.name
}

// ***

// check listener

type checkEventsLogger struct{}

func (l *checkEventsLogger) OnCheckRegistered(name string, res gosundheit.Result) {
	fmt.Println(fmt.Sprintf("[INFO] Check %q registered with initial result: %v", name, res))
}

func (l *checkEventsLogger) OnCheckStarted(name string) {
	fmt.Println(fmt.Sprintf("[INFO] Check %q started...", name))
}

func (l *checkEventsLogger) OnCheckCompleted(name string, res gosundheit.Result) {
	fmt.Println(fmt.Sprintf("[INFO] Check %q completed with result: %v", name, res))
}

// health listener

type checkHealthLogger struct{}

func (l *checkHealthLogger) OnResultsUpdated(results map[string]gosundheit.Result) {
	fmt.Println(fmt.Sprintf("[INFO] There are %d results, general health is %t\n", len(results), allHealthy(results)))
}

// INFO: duplicates of gosundheit/utils.go/allHealthy(..)
func allHealthy(results map[string]gosundheit.Result) (healthy bool) {
	for _, v := range results {
		if !v.IsHealthy() {
			return false
		}
	}

	return true
}

// ***

// api handler

type Api struct {
	Msg string `json:"message"`
}

func apiHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := json.NewEncoder(writer).Encode(&Api{Msg: "Hello, world!"})
		if err != nil {
			fmt.Println(fmt.Sprintf("[ERROR] Encoding api response failed: %s", err.Error()))
		}
	}
}
