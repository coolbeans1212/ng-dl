package filesystem

import (
	"fmt"
	"github.com/bogem/id3v2/v2"
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
	Rating      float32
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
	os.Rename(path, newPath)
}
