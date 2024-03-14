package targets

import (
	"context"
	"fmt"
	"goinload/internal/download"
	"goinload/internal/utils"

	"github.com/gocolly/colly"
)

// Scrape function accepts a link to scrape
func TargetBunkr(ctx context.Context, link string) error {
	// Create a new collector
	c := colly.NewCollector()

	// Set the user agent randomly
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
	})

	// Create a slice to store image URLs
	var galleryTitle string
	var imageLinks []string

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		if galleryTitle == "" {
			galleryTitle = e.Text
		}
	})

	// Find and collect product page links
	c.OnHTML(".grid-images a", func(e *colly.HTMLElement) {
		imagePageLink := e.Attr("href")

		c.Visit(imagePageLink)
	})

	// Find and collect product page links
	c.OnHTML(".container img", func(e *colly.HTMLElement) {
		imageLink := e.Attr("src")


		isImg := utils.IsImage(imageLink)

		if isImg {
			imageLinks = append(imageLinks, e.Request.AbsoluteURL(imageLink))
		}

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

	utils.EmitMsg(ctx, utils.LogEvent, "Start bulk Download", true)

	downloadPath, err := download.DownloadImages(ctx, imageLinks, galleryTitle, nil)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Images downloaded successfully.")
		utils.OpenFolder(downloadPath)
	}

	return nil
}
