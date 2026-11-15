package flags

import (
	"flag"
)

type CliFlags struct {
	DontFixMetadata bool
	DontRename      bool
	Resolution      string
	NoProgressBar   bool
}

func ParseFlags() CliFlags {
	var flags CliFlags
	flag.BoolVar(&flags.DontFixMetadata, "dont-fix-metadata", false, "Do not change the metadata of downloaded files.")
	flag.BoolVar(&flags.DontRename, "dont-rename", false, "Do not rename downloaded files, instead use the filename from the server.")
	flag.BoolVar(&flags.NoProgressBar, "no-progress", false, "Don't show a progress bar for downloads.")
	flag.StringVar(&flags.Resolution, "resolution", "highest", "[VIDEO] Get video of a specific resolution (e.g. 1080p).")
	flag.Parse()
	return flags
}
