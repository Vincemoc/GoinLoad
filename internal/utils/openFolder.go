package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenFolder(folderPath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", folderPath)
	case "linux": // Linux
		cmd = exec.Command("xdg-open", folderPath)
	case "windows": // Windows
		cmd = exec.Command("explorer", folderPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}
