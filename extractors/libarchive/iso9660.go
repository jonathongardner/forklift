package libarchive

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
)

const iso9660 = "application/octet-stream;iso=true"

func init() {
	isoDetector := helpers.MatchSigMultiOffsetFunc([]byte{0x43, 0x44, 0x30, 0x30, 0x31}, []int{0x8001, 0x8801, 0x9001})
	mimetype.Extend(isoDetector, iso9660, ".iso")
}
