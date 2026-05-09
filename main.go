package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/syumai/workers"
)

type ServerConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
	Auth string `json:"auth,omitempty"` // For Palworld basic auth if needed
}

type ServerStatus struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"` // "online" or "offline"
	Detail string `json:"detail"` // Additional info if available
	Ping   int64  `json:"ping"`   // ms
}

func fetchStatus(config ServerConfig) ServerStatus {
	start := time.Now()

	status := ServerStatus{
		Name: config.Name,
		Type: config.Type,
		Status: "offline",
		Detail: "",
	}

	req, err := http.NewRequest("GET", config.URL, nil)
	if err != nil {
		status.Detail = "Invalid URL"
		return status
	}

	if config.Auth != "" {
		req.Header.Add("Authorization", "Basic "+config.Auth)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	status.Ping = time.Since(start).Milliseconds()

	if err != nil {
		status.Detail = err.Error()
		return status
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status.Status = "online"

		if config.Type == "palworld" {
			body, _ := io.ReadAll(resp.Body)
			status.Detail = fmt.Sprintf("OK (Body: %d bytes)", len(body))
		} else {
			status.Detail = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
	} else {
		status.Detail = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return status
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	serversJSON := os.Getenv("SERVERS_JSON")
	var configs []ServerConfig

	if serversJSON != "" {
		if err := json.Unmarshal([]byte(serversJSON), &configs); err != nil {
			http.Error(w, "Failed to parse SERVERS_JSON", http.StatusInternalServerError)
			return
		}
	} else {
		configs = []ServerConfig{
			{Name: "Example", Type: "http", URL: "https://example.com"},
		}
	}

	var wg sync.WaitGroup
	statuses := make([]ServerStatus, len(configs))

	for i, config := range configs {
		wg.Add(1)
		go func(index int, conf ServerConfig) {
			defer wg.Done()
			statuses[index] = fetchStatus(conf)
		}(i, config)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(statuses)
}

func main() {
	http.HandleFunc("/api/status", apiHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
	})

	workers.Serve(nil)
}
