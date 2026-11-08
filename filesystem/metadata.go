package filesystem

import (
	"fmt"
	"github.com/bogem/id3v2/v2"
	"github.com/gcottom/mp4meta"
	"os"
	"path/filepath"
)

type Song struct {
	Title       string
	Artist      string
	Description string
	ArtistURL   string
	Year        string
	Genre       string
	Copyright   string
	DownloadURL string
	AlbumArtURL string
	Rating      float32
}

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

func starsFloatToPOPM(stars float32) (popm int) {
	if stars < 1.5 {
		return 1
	}
	if stars < 2.5 {
		return 64
	}
	if stars < 3.5 {
		return 128
	}
	if stars < 4.5 {
		return 196
	}
	if stars > 4.5 {
		return 255
	}
	return 0
}

func WriteMetadataToMP3File(filepath string, song Song) {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Println("Could not open downloaded MP3 file for writing. The MP3 file will not contain any metadata.")
		return
	}
	defer tag.Close()
	tag.SetVersion(3)
	encoding := id3v2.EncodingUTF8
	tag.SetTitle(song.Title)
	tag.SetAlbum(song.Title) // newgrounds doesn't have albums :(. set the album to the title so that iTunes doesn't show Unknown Album for everything
	tag.SetArtist(song.Artist)
	tag.SetYear(song.Year)
	tag.SetGenre(song.Genre)
	tag.AddTextFrame(tag.CommonID("Band/Orchestra/Accompaniment"), encoding, song.Artist)

	tag.DeleteFrames("POPM")
	starRatingBody := append([]byte("Windows Media Player 9 Series"), 0x00)
	starRatingBody = append(starRatingBody, byte(starsFloatToPOPM(song.Rating)))
	starRatingBody = append(starRatingBody, 0, 0, 0, 0)
	tag.AddFrame("POPM", id3v2.UnknownFrame{Body: starRatingBody}) // the star rating field

	tag.AddCommentFrame(id3v2.CommentFrame{Encoding: encoding, Language: "eng", Description: song.Description, Text: song.Description})
	tag.AddTextFrame("TCOP", encoding, song.Copyright)
	err = tag.Save()
	if err != nil {
		fmt.Println("Could not save metadata to file. The MP3 file will not contain any metadata.")
	}
	fmt.Println("Applied metadata to file.")
}

func RenameMP3File(path string, song Song) {
	directory := filepath.Dir(path)
	newName := song.Artist + " - " + song.Title + ".mp3"
	newPath := filepath.Join(directory, newName)
	fmt.Printf("Renamed file to %s\n", newName)
	err := os.Rename(path, newPath)
	if err != nil {
		fmt.Println("Could not rename the file.")
	}
}

func AttachAlbumArt(path string, albumArt []byte) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Printf("Could not open downloaded MP3 file for writing. The MP3 file will not contain any album art (%s).\n", err)
		return
	}
	defer tag.Close()
	tag.SetVersion(3)
	encoding := id3v2.EncodingUTF8
	tag.AddAttachedPicture(id3v2.PictureFrame{Encoding: encoding, MimeType: "image/png", PictureType: id3v2.PTFrontCover, Picture: albumArt})
	err = tag.Save()
	if err != nil {
		fmt.Println("Could not save album art to file. The MP3 file will not contain any album art.")
	}
	fmt.Println("Attached album art.")
}

func RenameFileInCurrentDirectoryAndPreserveExtension(filename string, renameto string) {
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get current directory.")
		return
	}
	fileExt := filepath.Ext(filename)
	err = os.Rename(filepath.Join(directory, filename), filepath.Join(directory, renameto)+fileExt)
	if err != nil {
		fmt.Printf("Could not rename the file: %s", err)
		return
	}
}

func WriteMetadataToMP4File(filepath string, video NGVideoList) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Could not open the file: %s\n", err)
		return
	}
	defer file.Close()
	mp4, err := mp4meta.ReadMP4(file)
	if err != nil {
		fmt.Printf("Could not attach metadata to the file: %s\n", err)
		return
	}
	mp4.SetAlbumArtist(video.Author)
	mp4.SetArtist(video.Author)
	mp4.SetTitle(video.Title)
	err = mp4.Save(file)
	if err != nil {
		fmt.Printf("Could not attach metadata to the file: %s\n", err)
		return
	}
}
