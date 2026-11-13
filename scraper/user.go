package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"net/url"
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

func DownloadAllAudioOfUser(user string) error {
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
	headers := http.Header{}
	headers.Set("X-Requested-With", "XMLHttpRequest")
	c.Request("GET", urlStr, nil, nil, headers)
	fmt.Println(audios)
	return err
}

func DownloadAllUserAssets(user string) {
	DownloadAllAudioOfUser(user)
}
