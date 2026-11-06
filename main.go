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
	if strings.Contains(url, "/audio/listen") {
		scraper.DownloadSongAndSaveWithMetadata(url)
	} else if strings.Contains(url, "newgrounds.com") {
		fmt.Println("Unimplemented.")
		os.Exit(1)
	} else {
		fmt.Println("Please provide a valid newgrounds.com link.")
		os.Exit(1)
	}

}
