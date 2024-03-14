package utils

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func EmitMsg(ctx context.Context, eventName EventType, message string, includeTimestamp bool) {
	var logMessage string

	if includeTimestamp {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logMessage = fmt.Sprintf("[%s] %s", timestamp, message)
	} else {
		logMessage = message
	}

	// Emit the log event with timestamped message
	runtime.EventsEmit(ctx, string(eventName), logMessage)
}

func IsImage(filename string) bool {
	// Convert the extension to lowercase for case-insensitive comparison
	extension := strings.ToLower(filepath.Ext(filename))

	// Check if the extension corresponds to an image format
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	for _, imgExt := range imageExtensions {
		if extension == imgExt {
			return true
		}
	}

	return false
}
