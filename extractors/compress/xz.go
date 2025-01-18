package compress

import (
	"io"

	"github.com/ulikunitz/xz"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	// log "github.com/sirupsen/logrus"
)

const Xz = "application/x-xz"

func decompressXz(r io.Reader) (io.Reader, error) {
	return xz.NewReader(r)
}

func ExtractXz(virtualFS *virtualfs.Fs) error {
	return helpers.ExtractCompression(virtualFS, decompressXz)
}
