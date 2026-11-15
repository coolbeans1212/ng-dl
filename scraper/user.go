package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/http"
	"net/url"
	"ng-dl/flags"
	"strconv"
	"strings"
	"sync"
	"time"
)

type UserPage struct {
	UsingElastic bool                `json:"using_elastic"`
	CalendarView bool                `json:"calendar_view"`
	Items        map[string][]string `json:"items"`
	Container    string              `json:"container"`
	LoadMore     string              `json:"load_more"`
	Url          string              `json:"url"`
	FreshURL     string              `json:"fresh_url"`
	Showcases    map[string]string   `json:"showcases"`
}

func GetSubdomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	parts := strings.Split(host, ".")

	if len(parts) <= 2 {
		return "", nil
	}

	return strings.Join(parts[:len(parts)-2], "."), nil
}

func GetAudioPageForUser(user string, page string) UserPage {
	var audios UserPage
	var err error = nil
	urlStr := "https://" + user + ".newgrounds.com/audio"
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &audios)
		if err != nil {
			fmt.Printf("Could not decode audio pages as JSON. Skipping audio...")
		}
	})
	u, err := url.Parse(urlStr)
	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()
	headers := http.Header{}
	headers.Set("X-Requested-With", "XMLHttpRequest")
	fmt.Printf("Requesting %s...\n", urlStr)
	c.Request("GET", u.String(), nil, nil, headers)
	return audios
}

func GetArtPageForUser(user string, page string) UserPage {
	var art UserPage
	var err error = nil
	urlStr := "https://" + user + ".newgrounds.com/art"
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &art)
		if err != nil {
			fmt.Printf("Could not decode art pages as JSON. Skipping art...\n")
		}
	})
	u, err := url.Parse(urlStr)
	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()
	headers := http.Header{}
	headers.Set("X-Requested-With", "XMLHttpRequest")
	fmt.Printf("Requesting %s...\n", urlStr)
	c.Request("GET", u.String(), nil, nil, headers)
	return art
}

func GetMoviePageForUser(user string, page string) UserPage {
	var movies UserPage
	var err error = nil
	urlStr := "https://" + user + ".newgrounds.com/movies"
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &movies)
		if err != nil {
			fmt.Printf("Could not decode movies pages as JSON. Skipping movies...\n")
		}
	})
	u, err := url.Parse(urlStr)
	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()
	headers := http.Header{}
	headers.Set("X-Requested-With", "XMLHttpRequest")
	fmt.Printf("Requesting %s...\n", urlStr)
	c.Request("GET", u.String(), nil, nil, headers)
	return movies
}

func DownloadAllAudioOfUser(user string, cliFlags flags.CliFlags) error {
	var err error = nil
	var wg sync.WaitGroup
	page := "1"
	for {
		audios := GetAudioPageForUser(user, page)
		fmt.Printf("Downloading page %s...\n", page)
		var doc *goquery.Document
		for i := range audios.Items {
			doc, err = goquery.NewDocumentFromReader(strings.NewReader(strings.Join(audios.Items[i], "")))
			doc.Find("a.item-audiosubmission").Each(func(i int, s *goquery.Selection) {
				url, exists := s.Attr("href")
				if !exists {
					return
				}
				if !cliFlags.NoParallel {
					wg.Add(1)
					go func(url string) {
						defer wg.Done()
						DownloadSongAndSaveWithMetadata(url, cliFlags)
					}(url)
					time.Sleep(time.Duration(float64(time.Second) * cliFlags.WaitBetweenParallelDownloads))
				} else {
					DownloadSongAndSaveWithMetadata(url, cliFlags)
				}
			})
		}
		if audios.LoadMore == "" {
			break
		} else {
			pageInt, _ := strconv.Atoi(page)
			page = strconv.Itoa(pageInt + 1)
		}
	}
	wg.Wait()
	return err
}

func DownloadAllArtOfUser(user string, cliFlags flags.CliFlags) error {
	var err error = nil
	var wg sync.WaitGroup
	page := "1"
	for {
		audios := GetArtPageForUser(user, page)
		fmt.Printf("Downloading page %s...\n", page)
		var doc *goquery.Document
		for i := range audios.Items {
			doc, err = goquery.NewDocumentFromReader(strings.NewReader(strings.Join(audios.Items[i], "")))
			doc.Find("a.item-portalitem-art-medium").Each(func(i int, s *goquery.Selection) {
				url, exists := s.Attr("href")
				if !exists {
					return
				}
				if !cliFlags.NoParallel {
					wg.Add(1)
					go func(url string) {
						defer wg.Done()
						DownloadArtAndSaveWithMetadata(url, cliFlags)
					}(url)
					time.Sleep(time.Duration(float64(time.Second) * cliFlags.WaitBetweenParallelDownloads))
				} else {
					DownloadArtAndSaveWithMetadata(url, cliFlags)
				}
			})
		}
		if audios.LoadMore == "" {
			break
		} else {
			pageInt, _ := strconv.Atoi(page)
			page = strconv.Itoa(pageInt + 1)
		}
	}
	wg.Wait()
	return err
}

func DownloadAllMoviesOfUser(user string, cliFlags flags.CliFlags) error {
	var err error = nil
	var wg sync.WaitGroup
	page := "1"
	for {
		audios := GetMoviePageForUser(user, page)
		fmt.Printf("Downloading page %s...\n", page)
		var doc *goquery.Document
		for i := range audios.Items {
			doc, err = goquery.NewDocumentFromReader(strings.NewReader(strings.Join(audios.Items[i], "")))
			doc.Find("a.inline-card-portalsubmission").Each(func(i int, s *goquery.Selection) {
				url, exists := s.Attr("href")
				if !exists {
					return
				}
				if !cliFlags.NoParallel {
					wg.Add(1)
					go func(url string) {
						defer wg.Done()
						DownloadVideo(url, cliFlags)
					}(url)
					time.Sleep(time.Duration(float64(time.Second) * cliFlags.WaitBetweenParallelDownloads))
				} else {
					DownloadVideo(url, cliFlags)
				}
			})
		}
		if audios.LoadMore == "" {
			break
		} else {
			pageInt, _ := strconv.Atoi(page)
			page = strconv.Itoa(pageInt + 1)
		}
	}
	wg.Wait()
	return err
}

func DownloadAllUserAssets(user string, cliFlags flags.CliFlags) {
	DownloadAllAudioOfUser(user, cliFlags)
	DownloadAllArtOfUser(user, cliFlags)
	DownloadAllMoviesOfUser(user, cliFlags)
}
