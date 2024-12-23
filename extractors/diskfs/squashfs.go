package diskfs

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
)

const SquashFS = "application/x-squashfs-image"

func init() {
	squashFSDetector := helpers.MatchSigFunc([]byte{0x68, 0x73, 0x71, 0x73}, 0)
	mimetype.Extend(squashFSDetector, SquashFS, ".squashfs")
}
