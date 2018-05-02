package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	max := 40
	min := 0
	var totalSamples int
	var sampleInterval int
	var minStep int
	var maxStep int
	var rnum int

	// Configure sample interval and totalSamples
	if duration == "hr" {
		// One sample every 10sec equates to 360 samples in an hour
		totalSamples = 360
		sampleInterval = 10
	} else {
		// One sample every 10min equates to 1008 samples in a week
		totalSamples = 1008
		sampleInterval = 10
	}

	if upward {
		if duration == "hr" {
			minStep = 2
			maxStep = 4
		} else {
			minStep = 3
			maxStep = 5
		}
	} else {
		minStep = 10
		maxStep = 10
	}

	var usageData []UsageSlice

	for sampleNum := 0; sampleNum < totalSamples; sampleNum += sampleInterval {
		// generate random number between max/min
		rnum = randomInt(min, max)

		usageData = append(usageData, UsageSlice{sampleNum, rnum})

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
	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	router := mux.NewRouter()
	router.HandleFunc("/usage", GetUsageHandler).Methods("GET").Queries("duration", "{duration}", "device", "{device}")

	fmt.Printf("Starting server on PORT %s...", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT),
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET"}),
			handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}
