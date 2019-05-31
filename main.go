package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	logger  = log.New(os.Stdout, "[BUTTONS] ", 0)
	secret  = mustEnvVar("secret", "")
	project = mustEnvVar("project", "")
	topic   = mustEnvVar("topic", "")
	que     *queue
)

func main() {

	q, err := newQueue(context.Background(), project, topic)
	if err != nil {
		logger.Fatal(err)
	}
	que = q

	http.HandleFunc("/", requestHandler)
	port := fmt.Sprintf(":%s", mustEnvVar("PORT", "8080"))
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal(err)
	}

}

func mustEnvVar(key, fallbackValue string) string {

	if val, ok := os.LookupEnv(key); ok {
		logger.Printf("%s: %s", key, val)
		return strings.TrimSpace(val)
	}

	if fallbackValue == "" {
		logger.Fatalf("Required envvar not set: %s", key)
	}

	logger.Printf("%s: %s (not set, using default)", key, fallbackValue)
	return fallbackValue
}
