package filesystem

import (
	"bytes"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"net/url"
	"ng-dl/console"
	"os"
	"path"
	"path/filepath"
)

func DownloadFileAndSave(urlStr string) string {
	var filename string
	progressid := rand.Int()
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	total := resp.ContentLength
	var downloaded int64
	var data []byte
	buffer := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			downloaded += int64(n)
			percent := 0
			if total > 0 {
				percent = int(float64(downloaded) / float64(total) * 100)
				if percent > 100 {
					percent = 100
				}
			}
			console.WriteProgress(progressid, "Downloading "+urlStr, percent)
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}
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
	err = os.WriteFile(out, data, 0644)
	if err != nil {
		fmt.Println("Could not save the file.")
	}
	console.ClearProgress(progressid)
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
