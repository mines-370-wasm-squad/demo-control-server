package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

type setStatusHandler struct {
	status string
}

func (h *setStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := writeStatus(h.status)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status, err := readStatus()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, strings.NewReader(status))
		if err != nil {
			log.Println(err)
		}
	})

	http.Handle("/enable", &setStatusHandler{"enabled"})
	http.Handle("/disable", &setStatusHandler{"disabled"})
	http.Handle("/stop", &setStatusHandler{"stopped"})

	http.ListenAndServe(":64684", nil)
}
