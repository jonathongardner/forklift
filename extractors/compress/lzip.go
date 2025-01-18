package compress

import (
	"io"

	"github.com/gabriel-vasile/mimetype"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	lzip "github.com/sorairolake/lzip-go"
	// log "github.com/sirupsen/logrus"
)

const Lzip = "application/lzip"

func init() {
	detector := helpers.MatchSigFunc([]byte{0x5D, 0x00, 0x00, 0x80, 0x00}, 0)
	mimetype.Extend(detector, Lzip, ".lzip")
}

func decompressLzip(r io.Reader) (io.Reader, error) {
	return lzip.NewReader(r)
}

func ExtractLzip(virtualFS *virtualfs.Fs) error {
	return helpers.ExtractCompression(virtualFS, decompressLzip)
}
