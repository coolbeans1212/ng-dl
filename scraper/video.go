package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"net/url"
	"ng-dl/filesystem"
	"os"
	"path"
	"strconv"
	"strings"
)

type NGVideoList struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Sources map[string][]struct {
		Type   string `json:"type"`
		Source string `json:"src"`
	} `json:"sources"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	} `json:"error"`
}

func GetVideoUrls(urlStr string) NGVideoList {
	var videoList NGVideoList
	id := path.Base(urlStr)
	source := "https://www.newgrounds.com/portal/video/" + id
	c := colly.NewCollector(colly.AllowedDomains("newgrounds.com", "www.newgrounds.com"))
	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &videoList)
		if err != nil {
			fmt.Printf("Could not decode video sources as JSON. Cannot continue.")
			os.Exit(1)
		}
	})
	headers := http.Header{}
	headers.Set("X-Requested-With", "XMLHttpRequest")
	fmt.Println("Finding video sources...")
	c.Request("GET", source, nil, nil, headers)
	return videoList
}

func DownloadHighestQualityVideo(urlStr string) {
	highestVerticalPx := 0
	highestResolution := "360p"
	videoList := GetVideoUrls(urlStr)
	for i := range videoList.Sources {
		var verticalPx int
		if i[len(i)-1:] == "k" {
			var err error
			verticalPx, err = strconv.Atoi(strings.TrimRight(i, "k"))
			if err != nil {
				fmt.Println("Could not get the video resolution.")
				continue
			}
			verticalPx *= 500
		} else {
			var err error
			verticalPx, err = strconv.Atoi(strings.TrimRight(i, "p"))
			if err != nil {
				fmt.Println("Could not get the video resolution.")
				continue
			}
		}
		if verticalPx > highestVerticalPx {
			highestVerticalPx = verticalPx
			highestResolution = i
		}
	}
	highestResolutionVideo := videoList.Sources[highestResolution][0].Source
	fmt.Printf("Downloading %s...\n", highestResolutionVideo)
	filesystem.DownloadFileAndSave(highestResolutionVideo)
	urlParsed, err := url.Parse(highestResolutionVideo)
	if err != nil {
		fmt.Println("Could not rename the file.")
	} else {
		filesystem.RenameFileInCurrentDirectoryAndPreserveExtension(path.Base(urlParsed.Path), fmt.Sprintf("%s - %s (%s)", videoList.Author, videoList.Title, strconv.Itoa(videoList.Id)))
	}
}
