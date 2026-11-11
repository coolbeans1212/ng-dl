package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SanitiseFileName(filename string) string {
	sanitiser := strings.NewReplacer(
		"<", "＜",
		">", "＞",
		":", "：",
		"\"", "＂",
		"/", "⧸",
		"\\", "⧹",
		"|", "｜",
		"?", "？",
		"*", "＊",
	)
	return sanitiser.Replace(filename)
}

func RenameFileInCurrentDirectoryAndPreserveExtension(filename string, renameto string) {
	filename = SanitiseFileName(filename)
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

func RenameMP3File(path string, song Song) {
	directory := filepath.Dir(path)
	newName := song.Artist + " - " + song.Title + ".mp3"
	newName = SanitiseFileName(newName)
	newPath := filepath.Join(directory, newName)
	fmt.Printf("Renamed file to %s\n", newName)
	err := os.Rename(path, newPath)
	if err != nil {
		fmt.Println("Could not rename the file.")
	}
}
