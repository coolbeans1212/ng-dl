package scraper

import (
	"net/url"
	"strings"
)

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

func DownloadAllUserAssets(urlStr string) {
	return
}
