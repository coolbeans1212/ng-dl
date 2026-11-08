package main

import (
	"flag"
	"fmt"
	"ng-dl/scraper"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: ng-dl [options] <url>")
		os.Exit(1)
	}
	url := args[0]
	switch {
	case strings.Contains(url, "/audio/listen"):
		scraper.DownloadSongAndSaveWithMetadata(url)
	case strings.Contains(url, "/playlist/") || strings.Contains(url, "/series/") || strings.Contains(url, "/collection/"):
		scraper.DownloadAllSongsInPlaylistWithMetadata(url)
		scraper.DownloadAllMoviesInPlaylistWithMetadata(url)
	case strings.Contains(url, "/portal/view/"):
		scraper.DownloadHighestQualityVideo(url)
	case strings.Contains(url, "newgrounds.com"):
		fmt.Println("Unimplemented.")
		os.Exit(1)
	default:
		fmt.Println("Please provide a valid newgrounds.com link.")
		os.Exit(1)
	}

}
