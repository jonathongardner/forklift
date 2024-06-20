package directory

import (
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	"github.com/jonathongardner/virtualfs/filetype"
	log "github.com/sirupsen/logrus"
)

func Add(add func(string, helpers.ExtratFunc)) {
	add(filetype.Dir.Mimetype, ExtractDir)
}

func ExtractDir(virtualFS *virtualfs.Fs) error {
	log.Debugf("Extracting directory %v", virtualFS.ErrorId())

	return nil
}
