package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
)

func GetMetadata(url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("newgrounds.com", "audio.ngfiles.com", "www.newgrounds.com"),
	)
	c.OnRequest(func (r *colly.Request) {
		fmt.Println("Fetching", r.URL.String())
	})
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if (e.Text == "NG Guard") {
			fmt.Println("Encountered NG Guard. Bypassing NG Guard is not implemented.")
			os.Exit(1)
		} else {
			title := e.Text
			fmt.Println("Song name is", title)
		}
	})

	err := c.Visit(url)
	if (err != nil) {
		fmt.Printf("Error %[1]s occured while trying to fetch %[2]s. Cannot continue.", err, url)
		os.Exit(1)
	}
}