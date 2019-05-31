package buttons

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	hookSecretEnvVarName      = "SEC"
	projectIDEnvVarName       = "PRJ"
	pubsubTopicNameEnvVarName = "TOP"
)

var (
	logger            = log.New(os.Stdout, "[BUTTONS] ", 0)
	once              sync.Once
	secret            string
	configInitializer = defaultConfigInitializer
	que               *queue
	errorJSON         = "{}"
)

func defaultConfigInitializer(ctx context.Context) error {

	logger.Print("Initializing configuration using: defaultConfigInitializer")

	secret = os.Getenv(hookSecretEnvVarName)
	if secret == "" {
		return fmt.Errorf("%s environment variable not set", hookSecretEnvVarName)
	}

	projectID := mustEnvVar(projectIDEnvVarName, "")
	topicName := mustEnvVar(pubsubTopicNameEnvVarName, "")

	q, err := newQueue(ctx, projectID, topicName)
	if err != nil {
		return err
	}
	que = q
	return nil
}

func mustEnvVar(key, fallbackValue string) string {

	if val, ok := os.LookupEnv(key); ok {
		logger.Printf("%s: %s", key, val)
		return val
	}

	if fallbackValue == "" {
		logger.Fatalf("Required envvar not set: %s", key)
	}

	logger.Printf("%s: %s (not set, using default)", key, fallbackValue)
	return fallbackValue
}
