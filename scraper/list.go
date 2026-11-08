package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strconv"
	"strings"
)

type VisualLinkSimple struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Url    string `json:"url"`
	Label  string `json:"label"`
	Status struct {
		Unpublished bool `json:"unpublished"`
	} `json:"status"`
	User struct {
		Id       int    `json:"id"`
		Username string `json:"user_name"`
		Url      string `json:"url"`
	} `json:"user"`
}

type VisualLinksResponse struct {
	Success  bool                                   `json:"success"`
	Partials map[string]map[string]string           `json:"partials"`
	Simples  map[string]map[string]VisualLinkSimple `json:"simples"`
	Wrapper  string                                 `json:"wrapper"`
}

func GetLinksFromPlaylist(urlStr string) map[string]map[string]VisualLinkSimple {
	var audioVisualLinkList []string
	c := colly.NewCollector(
		colly.AllowedDomains("newgrounds.com", "www.newgrounds.com"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Fetching", r.URL.String())
	})
	c.OnHTML(`ul#playlist_list`, func(e *colly.HTMLElement) {
		count := e.DOM.Children().Length()
		audioVisualLinkList = make([]string, count)
	})
	c.OnHTML(`ul#playlist_list li.visual-link-container`, func(e *colly.HTMLElement) {
		audioVisualLinkList[e.Index] = e.Attr("data-visual-link")
	})
	err := c.Visit(urlStr)
	if err != nil {
		fmt.Printf("Error %[1]s occured while trying to fetch %[2]s. Cannot continue.\n", err, urlStr)
		os.Exit(1)
	}
	// newgrounds doesn't put the playlist links in automatically, so we need to fetch them from /visual-links-fetch
	vlc := colly.NewCollector(colly.AllowedDomains("newgrounds.com", "www.newgrounds.com", "httpbin.org"))
	var visualLinksResponse VisualLinksResponse
	vlc.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	})
	vlc.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &visualLinksResponse)
		if err != nil {
			fmt.Printf("Could not decode the visual links response as JSON. Cannot continue.")
			os.Exit(1)
		}

	})
	var visualLinkForm map[string]string
	visualLinksStr := strings.Join(audioVisualLinkList, ",")
	visualLinksStr = fmt.Sprintf("[%s]", visualLinksStr)
	visualLinkForm = map[string]string{
		"component_params[context]": "playlist",
		"ids":                       visualLinksStr,
		"include_all_suitabilities": "0",
		"isAjaxRequest":             "1",
	}
	fmt.Println("Fetching visual links...")
	err = vlc.Post("https://www.newgrounds.com/visual-links-fetch", visualLinkForm)
	if err != nil {
		fmt.Printf("Error %[1]s occured while trying to fetch %[2]s. Cannot continue.\n", err, urlStr)
		os.Exit(1)
	}

	return visualLinksResponse.Simples
}

func DownloadAllSongsInPlaylistWithMetadata(urlStr string) {
	audioList := GetLinksFromPlaylist(urlStr)
	if len(audioList["3"]) == 0 {
		return
	}
	fmt.Printf("Downloading %s songs...\n", strconv.Itoa(len(audioList["3"])))
	for i := range audioList["3"] {
		DownloadSongAndSaveWithMetadata(audioList["3"][i].Url)
	}
}

func DownloadAllMoviesInPlaylistWithMetadata(urlStr string) {
	videoList := GetLinksFromPlaylist(urlStr)
	if len(videoList["1"]) == 0 {
		return
	}
	fmt.Printf("Downloading %s videos...\n", strconv.Itoa(len(videoList["1"])))
	for i := range videoList["1"] {
		DownloadHighestQualityVideo(videoList["1"][i].Url)
	}
}
