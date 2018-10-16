package starwarsapi

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"net/http"
	"os"
	"strconv"
)

var cb *gobreaker.CircuitBreaker

const cbRequestCountDefault = 3
const failureRatioDefault = 0.5

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	bootPhaseLogger := log.WithFields(log.Fields{
		"phase": "boot",
	})

	cbReqCount, err := strconv.Atoi(os.Getenv("CB_REQ_COUNT"))
	if err != nil {
		bootPhaseLogger.WithFields(log.Fields{
			"ENV VAR NAME":      "CB_REQ_COUNT",
			"ENV VAR VALUE":     cbReqCount,
			"APP DEFAULT VALUE": cbRequestCountDefault,
		}).Warn("Invalid circuit breaker request count. using default value")
		cbReqCount = cbRequestCountDefault
	}

	failureRatio, err := strconv.ParseFloat(os.Getenv("CB_REQ_FAIL_RATIO"), 64)
	if err != nil {
		bootPhaseLogger.WithFields(log.Fields{
			"ENV VAR NAME":      "CB_REQ_FAIL_RATIO",
			"ENV VAR VALUE":     failureRatio,
			"APP DEFAULT VALUE": failureRatioDefault,
		}).Warn("Invalid circuit breaker fail ratio. using default value")
		failureRatio = failureRatioDefault
	}

	var st gobreaker.Settings
	st.Name = "HTTP GET SWAPI"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= uint32(cbReqCount) && failureRatio >= failureRatio
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

// Wrapper para http.Get usando Circuit Breaker.
func GetPlanet(url string) (PlanetSearchResult, error) {
	result, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return PlanetSearchResult{}, err
		}
		defer resp.Body.Close()
		var p = &PlanetSearchResult{}
		json.NewDecoder(resp.Body).Decode(p)
		return p, nil
	})
	if err != nil {
		return PlanetSearchResult{}, err
	}

	return result.(PlanetSearchResult), nil
}
