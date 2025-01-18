package tar

import (
	"compress/gzip"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
)

const TarGz = "application/x-tar+gzip"

func init() {
	detector := helpers.MatchSigFunc([]byte{0x1f, 0x8b, 0x08, 0x00}, 0)
	mimetype.Extend(detector, TarGz, ".tar.gz")
}

func ExtractGzip(virtualFS *virtualfs.Fs) error {
	f, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer f.Close()

	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("could not open gzip reader (%v)", err)
	}
	return extract(virtualFS, gzipReader)
}
