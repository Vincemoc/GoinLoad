package utils

import (
	"context"
	"fmt"
	"io"
)

// ProgressWriter is an io.Writer implementation that tracks and prints progress.
type ProgressWriter struct {
	Writer  io.Writer
	Total   int64
	Written int64
	Filename string
	Ctx context.Context
}

// Write implements the io.Writer interface.
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	pw.Written += int64(n)

	progress := float64(pw.Written) / float64(pw.Total) * 100.0

	// Print progress, filename, and size to the console
	// log:= fmt.Sprintf("\rDownloading %s... %.2f%% (%s / %s)", pw.Filename, progress, formatBytes(pw.Written), formatBytes(pw.Total))
	log := fmt.Sprintf(`{"name":"%s","progress":"%.2f%%","written":"%s","total":"%s"}`, pw.Filename, progress, formatBytes(pw.Written), formatBytes(pw.Total))
	EmitMsg(pw.Ctx, ImageProgress, log, false)
	
	return n, err
}

// formatBytes formats bytes into a human-readable string.
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
