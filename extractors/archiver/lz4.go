package archiver

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
)

const xlz4 = "application/x-lz4"

func init() {
	isoDetector := helpers.MatchSigFunc([]byte{0x04, 0x22, 0x4D, 0x18}, 0)
	mimetype.Extend(isoDetector, xlz4, ".lz4")
}
