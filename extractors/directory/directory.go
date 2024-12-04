package directory

import (
	"github.com/jonathongardner/virtualfs"
	log "github.com/sirupsen/logrus"
)

func ExtractDir(virtualFS *virtualfs.Fs) error {
	log.Debugf("Extracting directory %v", virtualFS.ErrorId())

	return nil
}
