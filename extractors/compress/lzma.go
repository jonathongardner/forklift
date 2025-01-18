package compress

import (
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/ulikunitz/xz/lzma"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	// log "github.com/sirupsen/logrus"
)

const Lzma = "application/x-lzma"

func init() {
	detector := helpers.MatchSigFunc([]byte{0x5D, 0x00, 0x00, 0x80, 0x00}, 0)
	mimetype.Extend(detector, Lzma, ".lzma")
}

func decompressLzma(r io.Reader) (io.Reader, error) {
	return lzma.NewReader(r)
}

func ExtractLzma(virtualFS *virtualfs.Fs) error {
	return helpers.ExtractCompression(virtualFS, decompressLzma)
}
