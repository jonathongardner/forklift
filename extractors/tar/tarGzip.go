package tar

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
)

const TarGz = "application/x-tar+gzip"

func init() {
	tarGzDetector := helpers.MatchSigFunc([]byte{0x1f, 0x8b, 0x08, 0x00}, 0)
	mimetype.Extend(tarGzDetector, TarGz, ".tar.gz")
}
