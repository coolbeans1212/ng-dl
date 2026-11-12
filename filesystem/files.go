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

func CreateDirectoryInCurrentDirectory(dirname string) {
	dirname = SanitiseFileName(dirname)
	curdir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get current directory.")
		return
	}
	os.Mkdir(filepath.Join(curdir, dirname), 0644)
}

func RenameFileInDirectoryAndPreserveExtension(directory string, filename string, renameto string) {
	renameto = SanitiseFileName(renameto)
	fileExt := filepath.Ext(filename)
	err := os.Rename(filepath.Join(directory, filename), filepath.Join(directory, renameto)+fileExt)
	if err != nil {
		fmt.Printf("Could not rename the file: %s\n", err)
		return
	}
}

func RenameFileInCurrentDirectoryAndPreserveExtension(filename string, renameto string) {
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get current directory.")
		return
	}
	RenameFileInDirectoryAndPreserveExtension(directory, filename, renameto)
}

func RenameFileInRelativeDirectoryAndPreserveExtension(directory string, filename string, renameto string) {
	currdir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get current directory.")
		return
	}
	relativedir := filepath.Join(currdir, directory)
	RenameFileInDirectoryAndPreserveExtension(relativedir, filename, renameto)
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
