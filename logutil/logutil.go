package logutil

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

// Global variables for debug mode and logging cache
var (
	debugMode    bool                    // Determines if debug logs should be displayed
	loggedEvents = make(map[string]bool) // Cache to store logged events to prevent duplicates
	mutex        sync.Mutex
	// Mutex to synchronize access to loggedEvents
)

// Init initializes the logrus logger with JSON formatting and INFO level
func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// SetDebugMode enables or disables debug logging
func SetDebugMode(debug bool) {
	debugMode = debug
}

// GenerateCorrelationID generates a unique ID for tracking logs and events
func GenerateCorrelationID() string {
	return uuid.New().String()
}

// LogRelationalStart logs the start of an event if debug mode is enabled
func LogRelationalStart(correlationID, event string, additionalFields map[string]interface{}) *logrus.Entry {
	if !debugMode {
		return nil
	}

	fields := logrus.Fields{
		"event":         event,
		"correlationID": correlationID,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"status":        "started",
	}
	mergeFields(fields, additionalFields)

	entry := logrus.WithFields(fields)
	entry.Info("Event started")
	return entry
}

// LogRelationalEnd logs the end of an event if debug mode is enabled
func LogRelationalEnd(correlationID, event string, additionalFields map[string]interface{}) *logrus.Entry {
	if !debugMode {
		return nil
	}

	fields := logrus.Fields{
		"event":         event,
		"correlationID": correlationID,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"status":        "completed",
	}
	mergeFields(fields, additionalFields)

	entry := logrus.WithFields(fields)
	entry.Info("Event completed")
	return entry
}

// LogError logs an error event regardless of debug mode
func LogError(correlationID, event string, err error, additionalFields map[string]interface{}) {
	fields := logrus.Fields{
		"event":         event,
		"correlationID": correlationID,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"error":         err.Error(),
		"status":        "error",
	}
	mergeFields(fields, additionalFields)

	logrus.WithFields(fields).Error("Error occurred")
}

// LogOnce logs an event only once to prevent duplicate logs
func LogOnce(event string, err error, additionalFields map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	logKey := event
	if loggedEvents[logKey] {
		return
	}

	fields := logrus.Fields{
		"event":     event,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	mergeFields(fields, additionalFields)

	logrus.WithFields(fields).Info("Event logged once")
	loggedEvents[logKey] = true
}

// LogSuccess logs a successful event only once to prevent duplicate logs
func LogSuccess(event string, additionalFields map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	logKey := event
	if loggedEvents[logKey] {
		return
	}

	fields := logrus.Fields{
		"event":     event,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	mergeFields(fields, additionalFields)

	logrus.WithFields(fields).Info("Event logged successfully")
	loggedEvents[logKey] = true
}

// mergeFields merges additional fields into the base log fields
func mergeFields(baseFields logrus.Fields, additionalFields map[string]interface{}) {
	for k, v := range additionalFields {
		baseFields[k] = v
	}
}
