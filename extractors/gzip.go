package extractors

import (
	"compress/gzip"

	"github.com/jonathongardner/forklift/fs"
	// log "github.com/sirupsen/logrus"
)

const GzipMtype = "application/gzip"

func gzipExtract(entry *fs.Entry) ([]*fs.Entry, error) {
	file, err := entry.OpenTmp()
	if err != nil {
		return []*fs.Entry{}, err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return []*fs.Entry{}, err
	}

	newEntry, err := entry.ExtractedFile("", entry.Mode, gzipReader)
	if err != nil {
		return []*fs.Entry{}, err
	}

	return []*fs.Entry{newEntry}, nil
}

func init() {
	addExtractor(GzipMtype, gzipExtract)
}
