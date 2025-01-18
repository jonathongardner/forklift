package compress

import (
	"compress/bzip2"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	// log "github.com/sirupsen/logrus"
)

const Bzip2 = "application/x-bzip2"

func decompressBzip(r io.Reader) (io.Reader, error) {
	return bzip2.NewReader(r), nil
}

func ExtractBzip2(virtualFS *virtualfs.Fs) error {
	return helpers.ExtractCompression(virtualFS, decompressBzip)
}
