package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AudioPage struct {
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

func GetAudioPageForUser(user string, page string) AudioPage {
	var audios AudioPage
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

func DownloadAllAudioOfUser(user string) error {
	var err error = nil
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
				DownloadSongAndSaveWithMetadata(url)
			})
		}
		if audios.LoadMore == "" {
			break
		} else {
			pageInt, _ := strconv.Atoi(page)
			page = strconv.Itoa(pageInt + 1)
		}
	}
	return err
}

func DownloadAllUserAssets(user string) {
	DownloadAllAudioOfUser(user)
}
