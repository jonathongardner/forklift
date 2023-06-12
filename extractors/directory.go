package extractors

import (
	"os"

	"github.com/jonathongardner/forklift/fs"

	log "github.com/sirupsen/logrus"
)

const DirMtype = "directory/directory"

func dirExtract(entry *fs.Entry) ([]*fs.Entry, error) {
	dirs, err := os.ReadDir(entry.FullPath())
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		log.Infof("Test %v", dir)
	}

	return make([]*fs.Entry, 0), nil
}

func init() {
	addExtractor(DirMtype, dirExtract)
}
