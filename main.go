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
		scraper.GetMetadata(url)
	}

}
