package download

import (
	"context"
	"fmt"
	"goinload/internal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)


func downloadImage(ctx context.Context, url, imageFilepath string, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	defer response.Body.Close()

	file, err := os.Create(imageFilepath)
	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	// Get total file size
	size, _ := strconv.Atoi(response.Header.Get("Content-Length"))

	// Extract the filename from the filepath
	filename := filepath.Base(imageFilepath)

	// Create proxy writer to track download progress
	progressWriter := &utils.ProgressWriter{
		Writer: file,
		Total:  int64(size),
		Filename: filename,
		Ctx: ctx,
	}

	// Copy response body to file with progress tracking
	_, err = io.Copy(progressWriter, response.Body)
	if err != nil {
		errCh <- err
		return
	}
}

func DownloadImages(ctx context.Context, imageLinks []string, folderName string, filename *string) (string, error) {
	// Sanitize folderName to remove invalid characters for a directory name
	folderName = strings.TrimSpace(folderName)
	
	// Determine the appropriate download path based on the operating system
	var downloadPath string
	switch runtime.GOOS {
	case "windows":
		downloadPath = filepath.Join(os.Getenv("USERPROFILE"), "Downloads", folderName)
	case "darwin", "linux":
		downloadPath = filepath.Join(os.Getenv("HOME"), "Downloads", folderName)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Create the target folder if it doesn't exist
	if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
		return "", err
	}

	var wg sync.WaitGroup
	var errCh = make(chan error)

	files := make([]string, len(imageLinks))

	utils.EmitMsg(ctx, utils.InfoTask, "Downloading images...", false)

	for i, link := range imageLinks {
		wg.Add(1)
		
		// Check if filename is nil and create a default filename
		if filename == nil {
			defaultFilename := fmt.Sprintf("image_%d.jpg", i+1)
			filename = &defaultFilename
		}

		files[i] = *filename

		filepath := filepath.Join(downloadPath, *filename)
		go downloadImage(ctx, link, filepath, &wg, errCh)
	}

	wg.Wait()
	close(errCh)

	log1 := fmt.Sprintf("Images downloaded successfully! Check folder %v", downloadPath)
	utils.EmitMsg(ctx, utils.InfoTask, log1, false)

	// Check for errors from goroutines
	for err := range errCh {
		if err != nil {
			return downloadPath, err
		}
	}

	return downloadPath, nil
}

