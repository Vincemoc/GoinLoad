package targets

import (
	"context"
	"encoding/json"
	"fmt"
	"goinload/internal/download"
	"goinload/internal/utils"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

// Scrape function accepts a link to scrape
func TargetCyberdrop(ctx context.Context, link string) error {
	// Create a new collector
	c := colly.NewCollector()

	// Set the user agent randomly
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
	})

	// Create a slice to store image URLs
	var galleryTitle string
	var imageLinks []string

	// Find and collect product page links
	c.OnHTML(".image-container", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a.image", "href")
		name := e.ChildText("span.name")

		isImg := utils.IsImage(name)

		if isImg {
			imageLinks = append(imageLinks, e.Request.AbsoluteURL(link))
		}

	})

	c.OnHTML(".container", func(e *colly.HTMLElement) {
		galleryTitle = e.ChildText("h1#title")
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	err := c.Visit(link)
	if err != nil {
		return err
	}

	utils.EmitMsg(ctx, utils.GalleryName, fmt.Sprintf("%v  —  %v", galleryTitle, len(imageLinks)), false)

	utils.EmitMsg(ctx, utils.LogEvent, fmt.Sprintf("Found %d images", len(imageLinks)), true)

	utils.EmitMsg(ctx, utils.InfoTask, fmt.Sprintf("Extracting image url %v", len(imageLinks)), false)

	_, s_err := getImageSignedURL(ctx, galleryTitle, imageLinks)
	if s_err != nil {
		log.Fatal(s_err)
	}

	return nil
}

type JSONResponse struct {
	URL string `json:"url"`
	Name string `json:"name"`
}

func getImageSignedURL(ctx context.Context, galleryTitle string, links []string) ([]string, error) {
	baseAPIURL := "https://cyberdrop.me/api/f"
	var failedURLs []string
	var jsonResponseURLs []string
	var downloadPath string



	for i, link := range links {
		id := extractID(link)
		if id != "" {
			apiURL := fmt.Sprintf("%s/%s", baseAPIURL, id)

			log := fmt.Sprintf("Extracting images %d/%d", i+1, len(links))
			utils.EmitMsg(ctx, utils.InfoTask, log, false)

			response, err := makeAPIRequest(apiURL)
			if err != nil {
				fmt.Println("Error making request for", apiURL, ":", err)
				continue // Continue with the next iteration of the loop
			} 

			// Parse the JSON response to extract the "url" field
			var jsonResponse JSONResponse
			json_err := json.Unmarshal([]byte(response), &jsonResponse)
			if json_err != nil {
				fmt.Println("Error parsing JSON response:", json_err)
				failedURLs = append(failedURLs, apiURL)
				continue // Continue with the next iteration of the loop
			}

			log2 := fmt.Sprintf("Start Downloading %v", jsonResponse.Name)
			utils.EmitMsg(ctx, utils.LogEvent,log2 , true)
			path, err := download.DownloadImages(ctx, []string{jsonResponse.URL}, galleryTitle, &jsonResponse.Name)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Image downloaded successfully.")
				downloadPath = path
			}


			jsonResponseURLs = append(jsonResponseURLs, jsonResponse.URL)
		}
	}

	 utils.OpenFolder(downloadPath)
	 log2 := fmt.Sprintf("Failed Downloads %v", len(failedURLs))
	 utils.EmitMsg(ctx, utils.LogEvent,log2 , true)

	return jsonResponseURLs, nil
}

func extractID(link string) string {
	re := regexp.MustCompile(`\/f\/([a-zA-Z0-9]+)$`)
	matches := re.FindStringSubmatch(link)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}


func makeAPIRequest(apiURL string) (string, error) {
	// Introduce a delay to avoid rate limiting
	time.Sleep(1800 * time.Millisecond)
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}