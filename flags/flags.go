package flags

import (
	"flag"
)

type CliFlags struct {
	DontFixMetadata bool
	DontRename      bool
}

func ParseFlags() CliFlags {
	var flags CliFlags
	flag.BoolVar(&flags.DontFixMetadata, "dont-fix-metadata", false, "Do not change the metadata of downloaded files.")
	flag.BoolVar(&flags.DontRename, "dont-rename", false, "Do not rename downloaded files, instead use the filename from the server.")
	flag.Parse()
	return flags
}
