package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"ng-dl/filesystem"
	"os"
)

func GetArtMetadata(url string) filesystem.Artwork {
	artwork := filesystem.Artwork{}
	artwork.DownloadURLs = make(map[string]string)
	c := colly.NewCollector(
		colly.AllowedDomains("newgrounds.com", "art.ngfiles.com", "www.newgrounds.com"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Fetching", r.URL.String())
	})
	c.OnHTML("head", func(e *colly.HTMLElement) {
		artwork.Title = e.ChildAttr(`meta[name="title"]`, "content")
	})
	c.OnHTML("div.item-details-main h4 a", func(e *colly.HTMLElement) {
		artwork.Author = e.Text
	})
	c.OnHTML(`div[id^="art-container-"] div`, func(e *colly.HTMLElement) {
		regularSrc := e.ChildAttr("div div a img", "src")
		if regularSrc == "" {
			artwork.DownloadURLs[e.Attr("id")] = e.ChildAttr("div div a img", "data-smartload-src")
		} else {
			artwork.DownloadURLs[e.Attr("id")] = regularSrc
		}
	})
	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Error %[1]s occured while trying to fetch %[2]s. Cannot continue.\n", err, url)
		os.Exit(1)
	}
	return artwork
}

func DownloadArtAndSaveWithMetadata(url string) {
	metadata := GetArtMetadata(url)
	dirname := fmt.Sprintf("%s - %s", metadata.Author, metadata.Title)
	filesystem.CreateDirectoryInCurrentDirectory(dirname)
	for i := range metadata.DownloadURLs {
		filename := filesystem.DownloadFileAndSaveInRelativeDir(metadata.DownloadURLs[i], dirname)
		filesystem.RenameFileInDirectoryAndPreserveExtension(dirname, filename, i)
	}
}
