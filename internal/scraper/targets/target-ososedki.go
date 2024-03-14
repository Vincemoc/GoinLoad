package targets

import (
	"context"
	"fmt"
	"goinload/internal/download"
	"goinload/internal/utils"

	"github.com/gocolly/colly"
)

// Scrape function accepts a link to scrape
func TargetOsosedki(ctx context.Context, link string) error {
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
	c.OnHTML("#masonry .thumbs", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		// srcset := e.ChildAttr("a", "data-srcset")
		// name := e.ChildAttr("img.img-fluid", "alt")

		// https://ososedki.com/photos/-10000001_10004246

		isImg := utils.IsImage(link)

		if isImg {
			imageLinks = append(imageLinks, e.Request.AbsoluteURL(link))
		}

	})

	c.OnHTML("h1.text-white:nth-child(5)", func(e *colly.HTMLElement) {
		galleryTitle = e.Text
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
