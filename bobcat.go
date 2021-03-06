package main

import (
	json "encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Bobcat struct {
	Debug         bool
	ErrorLimit    int
	periodSeconds int
	address       string
	eventChannel  chan<- BobcatStatus
	errorCount    int
}

type BobcatStatus struct {
	Status           string `json:"status"`
	Gap              int64  `json:"gap"`
	MinerHeight      int64  `json:"miner_height"`
	BlockchainHeight int64  `json:"blockchain_height"`
	Epoch            int64  `json:"epoch"`
	LatencyMs        int64  `json:"latency_ms"`
	ErrorCount       int    `json:"error_count"`
	Valid            bool   `json:"valid"`
}

type bobcatStatusJson struct {
	Status           string `json:"status"`
	Gap              string `json:"gap"`
	MinerHeight      string `json:"miner_height"`
	BlockchainHeight string `json:"blockchain_height"`
	Epoch            string `json:"epoch"`
}

var HTTP_TIMEOUT_SECONDS = 30

func (bobcat Bobcat) Begin() {
	// SANITY CHECK
	if bobcat.address == "" {
		panic("Empty Bobcat address!")
	}
	if bobcat.ErrorLimit == 0 {
		bobcat.ErrorLimit = 1
	}

	bobcatUrl := url.URL{Scheme: "http", Host: bobcat.address, Path: "/status.json"}
	client := http.Client{
		Timeout: time.Duration(HTTP_TIMEOUT_SECONDS) * time.Second,
	}

	for {
		if bobcat.Debug {
			log.Printf("Fetching bobcat status from %s", bobcatUrl.String())
		}
		requestStart := time.Now()
		resp, err := client.Get(bobcatUrl.String())
		if err != nil {
			bobcat.errorCount++
			if bobcat.errorCount >= bobcat.ErrorLimit {
				bobcatStatus := BobcatStatus{
					Status:     "Error",
					LatencyMs:  int64(time.Since(requestStart) / time.Millisecond),
					ErrorCount: bobcat.errorCount,
					Valid:      false,
				}
				bobcat.eventChannel <- bobcatStatus
			}
			log.Printf("Error while fetching bobcat status: %s \n", err)
			time.Sleep(time.Duration(bobcat.periodSeconds) * time.Second)
			continue
		}
		requestEnd := time.Since(requestStart)
		bobcat.errorCount = 0
		if bobcat.Debug {
			log.Printf("Got response from bobcat! Latency: %v\n", requestEnd)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error while reading response from bobcat : %s \n", err)
			time.Sleep(time.Duration(bobcat.periodSeconds) * time.Second)
			continue
		}

		jsonResponse := &bobcatStatusJson{}
		if err := json.Unmarshal(body, &jsonResponse); err != nil {
			log.Printf("Error parsing JSON body from bobcat (\"%s\"): %s \n", string(body), err)
			time.Sleep(time.Duration(bobcat.periodSeconds) * time.Second)
			continue
		}

		// PARSE STRINGFUL JSON INTO STRUCT WITH INTS
		bobcatStatus := BobcatStatus{
			Status:    jsonResponse.Status,
			LatencyMs: int64(requestEnd / time.Millisecond),
			Valid:     true,
		}
		gap, err := strconv.ParseInt(jsonResponse.Gap, 10, 64)
		if err != nil {
			bobcatStatus.Valid = false
		} else {
			bobcatStatus.Gap = gap
		}

		minerHeight, err := strconv.ParseInt(jsonResponse.MinerHeight, 10, 64)
		if err != nil {
			bobcatStatus.Valid = false
		} else {
			bobcatStatus.MinerHeight = minerHeight
		}

		blockchainHeight, err := strconv.ParseInt(jsonResponse.BlockchainHeight, 10, 64)
		if err != nil {
			bobcatStatus.Valid = false
		} else {
			bobcatStatus.BlockchainHeight = blockchainHeight
		}

		epoch, err := strconv.ParseInt(jsonResponse.Epoch, 10, 64)
		if err != nil {
			bobcatStatus.Valid = false
		} else {
			bobcatStatus.Epoch = epoch
		}

		bobcat.eventChannel <- bobcatStatus

		time.Sleep(time.Duration(bobcat.periodSeconds) * time.Second)
	}
}
