/*
 * Copyright (c) 2024, Paul Gundarapu.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

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
