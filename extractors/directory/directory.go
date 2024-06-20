package directory

import (
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/forklift/filetype"
	"github.com/jonathongardner/forklift/fs"
	log "github.com/sirupsen/logrus"
)

func Add(add func(string, helpers.ExtratFunc)) {
	add(filetype.Dir.Mimetype, ExtractDir)
}

func ExtractDir(virtualFS *fs.Virtual) error {
	log.Debugf("Extracting directory %v", virtualFS.RootErrorID())

	return nil
}
