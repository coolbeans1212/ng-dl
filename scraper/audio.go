package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"ng-dl/filesystem"
	"os"
	"strconv"
	"strings"
)

func GetAudioMetadata(url string) filesystem.Song {
	var title, metaTitle, metaDescription, artist, authorUrl, year, genre, copyright, downloadLink, rating, albumArtURL string
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
		if e.Index != 0 {
			return
		}
		copyright = e.Text
		copyright = strings.TrimSpace(copyright)
	})
	c.OnHTML(`div#creative_commons div dl dt`, func(e *colly.HTMLElement) {
		if e.Index == 0 {
			copyright += " Creative Commons: "
		}
		copyright += strings.TrimRight(e.Text, ":") + " " //TrimRight because the string is often something like "Attribution:" so we need to remove the ":"
	})
	c.OnHTML("#score_number", func(e *colly.HTMLElement) {
		rating = e.Text
	})
	c.OnHTML("img.audio-cover-art", func(e *colly.HTMLElement) {
		albumArtURL = e.Attr("src")
	})
	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Error %[1]s occured while trying to fetch %[2]s. Cannot continue.\n", err, url)
		os.Exit(1)
	}
	fmt.Printf("Fetched metadata for %s\n", title)
	ratingFloat, err := strconv.ParseFloat(rating, 32)
	if err != nil {
		fmt.Println("Fetched rating is invalid. Defaulting to 4.")
		ratingFloat = 4.0
	}
	return filesystem.Song{Title: title, Description: metaDescription, Artist: artist, ArtistURL: authorUrl, Year: year, Genre: genre, Copyright: copyright, DownloadURL: downloadLink, Rating: float32(ratingFloat), AlbumArtURL: albumArtURL}
}

func DownloadSongAndSaveWithMetadata(url string) {
	metadata := GetAudioMetadata(url)
	filepath := filesystem.DownloadFileAndSave(metadata.DownloadURL)
	filesystem.WriteMetadataToMP3File(filepath, metadata)
	filesystem.AttachAlbumArt(filepath, filesystem.DownloadAndConvertWebpToPng(metadata.AlbumArtURL))

	filesystem.RenameMP3File(filepath, metadata)
}
