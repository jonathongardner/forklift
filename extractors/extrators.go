package extractors

import (
	"github.com/jonathongardner/forklift/extractors/directory"
	"github.com/jonathongardner/virtualfs/filetype"

	// "github.com/jonathongardner/forklift/extractors/gzip"
	"github.com/jonathongardner/forklift/extractors/diskfs"
	"github.com/jonathongardner/forklift/extractors/helpers"
	// log "github.com/sirupsen/logrus"
)

var Functions = make(map[string]helpers.ExtratFunc)
var Types []string

type extractor struct {
	mtype string
	ext   helpers.ExtratFunc
}

// extracts to folder
func addExtractor(ext extractor) {
	if _, ok := Functions[ext.mtype]; ok {
		panic("extractor already exists")
	}
	Functions[ext.mtype] = ext.ext
	Types = append(Types, ext.mtype)
}

func init() {
	toAdd := []extractor{
		//-----------------directory-----------------
		{filetype.Dir.Mimetype, directory.ExtractDir},
		//-----------------directory-----------------

		//-----------------diskfs-----------------
		{diskfs.SquashFS, diskfs.ExtractArchive},
		//-----------------diskfs-----------------

	}
	for _, t := range toAdd {
		addExtractor(t)
	}
}
