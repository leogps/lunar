package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

var logger *slog.Logger

// YeetLogHandler formats log entries according to your desired format
type YeetLogHandler struct {
	*slog.TextHandler
}

// Handle formats and outputs the log entry
func (h *YeetLogHandler) Handle(_ context.Context, r slog.Record) error {
	//// Format the timestamp
	//timestamp := r.Time.Format(time.RFC3339)
	//
	//// Format the log level
	//level := r.Level.String()

	// Format the log message
	msg := r.Message

	// Print the log in the desired format
	//_, err := fmt.Fprintf(os.Stdout, "%s %s - %s\n", timestamp, level, msg)
	_, err := fmt.Fprintf(os.Stdout, "%s\n", msg)
	if err != nil {
		return err
	}
	return nil
}

func InitLogger(level slog.Level) {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	})
	// Create the custom handler
	customHandler := &YeetLogHandler{
		textHandler,
	}

	logger = slog.New(customHandler)
}

func LogInfo(message string, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf(message, args...))
		return
	}
	fmt.Printf(fmt.Sprintf(message+"\n", args...))
}

func LogDebug(message string, args ...interface{}) {
	if logger != nil {
		logger.Debug(fmt.Sprintf(message, args...))
		return
	}
	fmt.Printf(fmt.Sprintf(message+"\n", args...))
}

func LogWarn(message string, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf(message, args...))
		return
	}
	fmt.Printf(fmt.Sprintf(message+"\n", args...))
}

func LogError(message string, err error, args ...interface{}) {
	var errorMessage string
	if err != nil {
		errorMessage = fmt.Sprintf(" error: %v", err)
	}
	effectiveMessage := fmt.Sprintf(message+errorMessage, args...)
	if logger != nil {
		logger.Error(effectiveMessage)
		return
	}
	fmt.Printf(fmt.Sprintf(effectiveMessage+"\n", args...))
}
