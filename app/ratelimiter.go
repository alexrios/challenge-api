package app

import (
	log "github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
	"net/http"
	"os"
	"strconv"
)

var rl ratelimit.Limiter;

const defaultRate = 100

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	bootPhaseLogger := log.WithFields(log.Fields{
		"phase": "boot",
	})

	rateLimitPerSecond := os.Getenv("API_RATE_LIMIT")
	rate, err := strconv.Atoi(rateLimitPerSecond)
	if err != nil {
		bootPhaseLogger.WithFields(log.Fields{
			"ENV VAR NAME":  "API_RATE_LIMIT",
			"ENV VAR VALUE": rate,
		}).Warn("Invalid Rate limit")
		rl = ratelimit.New(defaultRate)
	} else {
		rl = ratelimit.New(rate)
	}
}

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

func RateLimitedHandler(next httpHandlerFunc) httpHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		rl.Take()
		next(res, req)
	}
}
