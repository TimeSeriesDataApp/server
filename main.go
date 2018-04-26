package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github/gorilla/mux"
)

//******************************************************************************
// Cpu Usage
type CpuLoadSlice struct {
	Toffset int `json:"toffset,omitempty"`
	Load    int `json:"load,omitempty"`
}

func GetCpuUsageHandler(w http.ResponseWriter, req *http.Request) {
	duration := req.URL.Query().Get("duration")

	if duration != "wk" && duration != "hr" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := randomCpuUsageData(duration)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func randomCpuUsageData(duration string) []CpuLoadSlice {
	// rand.Intn(max - min) + min
	max := 100
	min := 0
	// go to 3600
	tsec_end := 3600
	tsec_interval := 5
	var rnum int

	var cpuLoad []CpuLoadSlice

	// generate a bunch of random numbers
	for tsec := 0; tsec < tsec_end; tsec += tsec_interval {
		// generate random number between max/min
		rnum = randomInt(min, max)

		cpuLoad = append(cpuLoad, CpuLoadSlice{tsec, rnum})

		// reassign max
		if (rnum + 10) > 100 {
			max = 100
		} else {
			max = rnum + 10
		}

		// reassign min
		if (rnum - 10) < 0 {
			min = 0
		} else {
			min = rnum - 10
		}
	}

	return cpuLoad
}

//******************************************************************************
// Disk Usage
func GetDiskUsageHandler(w http.ResponseWriter, req *http.Request) {

}

//******************************************************************************
// Memory Usage
func GetMemoryUsageHandler(w http.ResponseWriter, req *http.Request) {

}

//******************************************************************************
// Network Usage
func GetNetworkSpeedHandler(w http.ResponseWriter, req *http.Request) {

}

func randomInt(min int, max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	return r.Intn(max-min) + min
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/usage/cpu", GetCpuUsageHandler).Methods("GET").Queries("duration", "{duration}")
	router.HandleFunc("/usage/disk", GetDiskUsageHandler).Methods("GET").Queries("duration", "{duration}")
	router.HandleFunc("/usage/memory", GetMemoryUsageHandler).Methods("GET").Queries("duration", "{duration}")
	router.HandleFunc("/usage/network", GetNetworkSpeedHandler).Methods("GET").Queries("duration", "{duration}")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
