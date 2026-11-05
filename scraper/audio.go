package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"ng-dl/filesystem"
	"os"
	"strings"
)

func GetMetadata(url string) {
	var title, metaTitle, metaDescription, artist, authorUrl, year, genre, copyright, downloadLink string
	c := colly.NewCollector(
		colly.AllowedDomains("newgrounds.com", "audio.ngfiles.com", "www.newgrounds.com"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Fetching", r.URL.String())
	})
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if e.Text == "NG Guard" {
			fmt.Println("Encountered NG Guard. Bypassing NG Guard is not implemented.")
			os.Exit(1)
		} else {
			title = e.Text
			fmt.Println("Song name is", title)
		}
	})
	c.OnHTML("head", func(e *colly.HTMLElement) {
		metaTitle = e.ChildAttr(`meta[name="title"]`, "content")
		metaDescription = e.ChildAttr(`meta[property="og:description"]`, "content")
		downloadLink = e.ChildAttr(`meta[property="og:audio"]`, "content")
		if metaTitle != title && metaTitle != "" {
			title = metaTitle
			fmt.Println("Never mind, song name is", title)
		}
		fmt.Println("Description:", metaDescription)
	})
	c.OnHTML("div.item-details-main h4 a", func(e *colly.HTMLElement) {
		artist = e.Text
		authorUrl = e.Attr("href")
	})
	c.OnHTML("div#sidestats dl.sidestats dd.multivalue span.value", func(e *colly.HTMLElement) {
		if e.Index != 0 {
			return
		}
		year = e.Text[len(e.Text)-4:]
	})
	c.OnHTML(`[id*="genre-view"] dl dd a`, func(e *colly.HTMLElement) {
		genre = e.Text
	})
	c.OnHTML(`div#creative_commons div p`, func(e *colly.HTMLElement) {
		copyright = e.Text
		copyright = strings.TrimSpace(copyright)
	})
	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Error %[1]s occured while trying to fetch %[2]s. Cannot continue.", err, url)
		os.Exit(1)
	}
	fmt.Println(title, metaDescription, artist, authorUrl, year, genre, copyright, "goog", downloadLink)
	filesystem.DownloadFileAndSave(downloadLink)
}
