package filesystem

import (
	"bytes"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func DownloadFileAndSave(urlStr string) string {
	var filename string
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if filename == "" {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			panic(err)
		}
		filename = path.Base(parsedURL.Path)
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	out := filepath.Join(dir, filename)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(out, data, 0644)
	if err != nil {
		panic(err)
	}
	return out

}

func DownloadAndConvertWebpToPng(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not fetch the album art.")
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Could not decode the album art.")
	}
	var buffer bytes.Buffer
	err = png.Encode(&buffer, img)
	return buffer.Bytes()
}
