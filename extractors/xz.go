package extractors

import (
	"github.com/jonathongardner/forklift/fs"
	"github.com/ulikunitz/xz"
	// log "github.com/sirupsen/logrus"
)

const XZMtype = "application/x-xz"

func xzExtract(entry *fs.Entry) ([]*fs.Entry, error) {
	file, err := entry.OpenTmp()
	if err != nil {
		return []*fs.Entry{}, err
	}
	defer file.Close()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return []*fs.Entry{}, err
	}

	newEntry, err := entry.ExtractedFile("", entry.Mode, xzReader)
	if err != nil {
		return []*fs.Entry{}, err
	}

	return []*fs.Entry{newEntry}, nil
}

func init() {
	addExtractor(XZMtype, xzExtract)
}
