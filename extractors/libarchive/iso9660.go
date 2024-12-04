package libarchive

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
)

const Iso9660 = "application/x-iso9660-image"

func init() {
	isoDetector := helpers.MatchSigMultiOffsetFunc([]byte{0x43, 0x44, 0x30, 0x30, 0x31}, []int{0x8001, 0x8801, 0x9001})
	mimetype.Extend(isoDetector, Iso9660, ".iso")
}
