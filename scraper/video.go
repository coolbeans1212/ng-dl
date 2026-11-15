package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"net/url"
	"ng-dl/filesystem"
	"ng-dl/flags"
	"os"
	"path"
	"strconv"
	"strings"
)

func GetVideoUrls(urlStr string) filesystem.NGVideoList {
	var videoList filesystem.NGVideoList
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

func GetHighestResolution(videoList filesystem.NGVideoList) string {
	highestVerticalPx := 0
	highestResolution := "360p"
	if videoList.Error.Code != 0 {
		fmt.Printf("Cannot download video: are you sure this is a video? %s", videoList.Error.Message)
		return ""
	}
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
	return highestResolution
}

func DownloadVideo(urlStr string, cliFlags flags.CliFlags) {
	videoList := GetVideoUrls(urlStr)
	var videoResolution string
	if cliFlags.Resolution == "highest" {
		videoResolution = GetHighestResolution(videoList)
	} else {
		videoResolution = cliFlags.Resolution
	}
	var video string
	sources, exists := videoList.Sources[videoResolution]
	if exists && len(sources) > 0 {
		video = videoList.Sources[videoResolution][0].Source
	} else {
		videoResolution = GetHighestResolution(videoList)
		video = videoList.Sources[videoResolution][0].Source
		fmt.Printf("Could not find resolution %[1]s, defaulting to %[2]s.\n", cliFlags.Resolution, videoResolution)
	}

	fmt.Printf("Downloading %s...\n", video)
	filesystem.DownloadFileAndSave(video, cliFlags)
	urlParsed, parseerr := url.Parse(video)
	if parseerr != nil {
		fmt.Printf("Could not get the file name: %s\n", parseerr)
		return
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get the current directory.")
	} else {
		if !cliFlags.DontFixMetadata {
			fmt.Println("Writing metadata...")
			filesystem.WriteMetadataToMP4File(path.Join(dir, path.Base(urlParsed.Path)), videoList)
		}
	}
	if !cliFlags.DontRename {
		newName := fmt.Sprintf("%s - %s (%s)", videoList.Author, videoList.Title, strconv.Itoa(videoList.Id))
		fmt.Printf("Renaming file to %s.mp4...\n", newName)
		filesystem.RenameFileInCurrentDirectoryAndPreserveExtension(path.Base(urlParsed.Path), newName)
	}
}

func DownloadFeaturedMovies(cliFlags flags.CliFlags) {
	videoList := make(map[int]string)
	c := colly.NewCollector(colly.AllowedDomains("www.newgrounds.com", "newgrounds.com"))
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if e.Text == "NG Guard" {
			fmt.Println("Encountered NG Guard. Bypassing NG Guard is not implemented.")
			os.Exit(1)
		}
	})
	c.OnHTML("div.featured-submission-grid--movies a.item-portalsubmission-featured", func(e *colly.HTMLElement) {
		videoList[e.Index] = e.Attr("href")
	})
	c.Visit("https://www.newgrounds.com/")
	for i := range videoList {
		DownloadVideo(videoList[i], cliFlags)
	}
}
