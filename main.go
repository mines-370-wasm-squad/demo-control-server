package main

import (
	"log"
	"net/http"
	"os"
)

const CONTROL_FILE = "/tmp/gpio-demo-control/status"

func readStatus() (string, error) {
	data, err := os.ReadFile(CONTROL_FILE)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func writeStatus(status string) error {
	return os.WriteFile(CONTROL_FILE, []byte(status), 0o644)
}

var nextState = map[string]string{
	"ready":   "enqueued",
	"running": "dequeued",
}

func main() {
	http.HandleFunc("/advance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status, err := readStatus()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next, ok := nextState[status]
		if !ok {
			log.Println("advance called when status is", status)
			http.Error(w, status, http.StatusServiceUnavailable)
			return
		}

		writeStatus(next)
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":64684", nil)
}
