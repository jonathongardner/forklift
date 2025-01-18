package compress

import (
	"compress/gzip"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	// log "github.com/sirupsen/logrus"
)

const Gzip = "application/gzip"

func decompressGzip(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}

func ExtractGzip(virtualFS *virtualfs.Fs) error {
	return helpers.ExtractCompression(virtualFS, decompressGzip)
}
