package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// UsageSlice : Time slice struct to hold generated usage data
type UsageSlice struct {
	Toffset int `json:"toffset"`
	Usage   int `json:"usage"`
}

// GetUsageHandler : Handler for the /usage endpoint
func GetUsageHandler(w http.ResponseWriter, req *http.Request) {
	duration := req.URL.Query().Get("duration")
	if duration != "wk" && duration != "hr" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid duration: %v\n", duration)
		return
	}

	// usageMap holds all generated data for all devices
	var usageMap = make(map[string][]UsageSlice)

	devices := strings.Split(req.URL.Query().Get("device"), ",")
	for i := range devices {
		if _, ok := usageMap[devices[i]]; ok {
			// Duplicate query in query string
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Duplicate device: %v\n", devices[i])
			return
		}

		upward := false
		if devices[i] == "disk" {
			// Disk usage will gradually trend upward
			upward = true
		}

		switch devices[i] {
		case "cpu", "disk", "memory", "network":
			usageMap[devices[i]] = randomUsageData(duration, upward)
		default:
			// Unsupported/unknown device
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown device: %v\n", devices[i])
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usageMap)
}

func randomUsageData(duration string, upward bool) []UsageSlice {
	max := 30
	min := 0
	var tsecEnd int
	var tsecInterval int
	var minStep int
	var maxStep int
	var rnum int

	// Configure sample interval and time end
	if duration == "hr" {
		// number of seconds in an hour
		tsecEnd = 3600
		tsecInterval = 5
	} else {
		// number of seconds in a week
		tsecEnd = 604800
		tsecInterval = 240
	}

	if upward {
		minStep = 8
		maxStep = 15
	} else {
		minStep = 10
		maxStep = 10
	}

	var usageData []UsageSlice

	// Generate a bunch of random numbers
	for tsec := 0; tsec < tsecEnd; tsec += tsecInterval {
		// generate random number between max/min
		rnum = randomInt(min, max)

		usageData = append(usageData, UsageSlice{tsec, rnum})

		// Reassign max
		if (rnum + maxStep) > 100 {
			max = 100
		} else {
			max = rnum + maxStep
		}

		// Reassign min
		if (rnum - minStep) < 0 {
			min = 0
		} else {
			min = rnum - minStep
		}
	}

	return usageData
}

func randomInt(min int, max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	return r.Intn(max-min) + min
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/usage", GetUsageHandler).Methods("GET").Queries("duration", "{duration}", "device", "{device}")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
