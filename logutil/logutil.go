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

// Struct to hold additional fields
type LogFields struct {
	Event         string                 `json:"event,omitempty"`
	CorrelationID string                 `json:"correlationID,omitempty"`
	Timestamp     string                 `json:"timestamp,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Additional    map[string]interface{} `json:"additional,omitempty"`
}

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

// LogRelationalStart logs the start of an event if debug mode is enabled using map[string]interface{}
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

// LogRelationalEnd logs the end of an event if debug mode is enabled using map[string]interface{}
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

// LogError logs an error event regardless of debug mode using map[string]interface{}
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

// LogOnce logs an event only once to prevent duplicate logs using map[string]interface{}
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

// LogSuccess logs a successful event only once to prevent duplicate logs using map[string]interface{}
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

// LogRelationalStartNew logs the start of an event if debug mode is enabled using struct
func LogRelationalStartNew(correlationID, event string, fields LogFields) *logrus.Entry {
	if !debugMode {
		return nil
	}

	fields.Event = event
	fields.CorrelationID = correlationID
	fields.Timestamp = time.Now().UTC().Format(time.RFC3339)
	fields.Status = "started"

	entry := logrus.WithFields(logrus.Fields{
		"event":         fields.Event,
		"correlationID": fields.CorrelationID,
		"timestamp":     fields.Timestamp,
		"status":        fields.Status,
	})
	mergeFieldsNew(entry, fields.Additional)

	entry.Info("Event started")
	return entry
}

// LogRelationalEndNew logs the end of an event if debug mode is enabled using struct
func LogRelationalEndNew(correlationID, event string, fields LogFields) *logrus.Entry {
	if !debugMode {
		return nil
	}

	fields.Event = event
	fields.CorrelationID = correlationID
	fields.Timestamp = time.Now().UTC().Format(time.RFC3339)
	fields.Status = "completed"

	entry := logrus.WithFields(logrus.Fields{
		"event":         fields.Event,
		"correlationID": fields.CorrelationID,
		"timestamp":     fields.Timestamp,
		"status":        fields.Status,
	})
	mergeFieldsNew(entry, fields.Additional)

	entry.Info("Event completed")
	return entry
}

// LogErrorNew logs an error event regardless of debug mode using struct
func LogErrorNew(correlationID, event string, err error, fields LogFields) {
	fields.Event = event
	fields.CorrelationID = correlationID
	fields.Timestamp = time.Now().UTC().Format(time.RFC3339)
	fields.Error = err.Error()
	fields.Status = "error"

	entry := logrus.WithFields(logrus.Fields{
		"event":         fields.Event,
		"correlationID": fields.CorrelationID,
		"timestamp":     fields.Timestamp,
		"error":         fields.Error,
		"status":        fields.Status,
	})
	mergeFieldsNew(entry, fields.Additional)

	entry.Error("Error occurred")
}

// LogOnceNew logs an event only once to prevent duplicate logs using struct
func LogOnceNew(event string, err error, fields LogFields) {
	mutex.Lock()
	defer mutex.Unlock()

	logKey := event
	if loggedEvents[logKey] {
		return
	}

	fields.Event = event
	fields.Timestamp = time.Now().UTC().Format(time.RFC3339)
	if err != nil {
		fields.Error = err.Error()
	}

	entry := logrus.WithFields(logrus.Fields{
		"event":     fields.Event,
		"timestamp": fields.Timestamp,
	})
	mergeFieldsNew(entry, fields.Additional)

	entry.Info("Event logged once")
	loggedEvents[logKey] = true
}

// LogSuccessNew logs a successful event only once to prevent duplicate logs using struct
func LogSuccessNew(event string, fields LogFields) {
	mutex.Lock()
	defer mutex.Unlock()

	logKey := event
	if loggedEvents[logKey] {
		return
	}

	fields.Event = event
	fields.Timestamp = time.Now().UTC().Format(time.RFC3339)

	entry := logrus.WithFields(logrus.Fields{
		"event":     fields.Event,
		"timestamp": fields.Timestamp,
	})
	mergeFieldsNew(entry, fields.Additional)

	entry.Info("Event logged successfully")
	loggedEvents[logKey] = true
}

// mergeFields merges additional fields into the base log fields (for map[string]interface{})
func mergeFields(baseFields logrus.Fields, additionalFields map[string]interface{}) {
	for k, v := range additionalFields {
		baseFields[k] = v
	}
}

// mergeFieldsNew merges additional fields into the base log fields (for LogFields struct)
func mergeFieldsNew(entry *logrus.Entry, additionalFields map[string]interface{}) {
	if additionalFields != nil {
		for k, v := range additionalFields {
			entry = entry.WithField(k, v)
		}
	}
}
