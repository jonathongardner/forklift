//go:build !libarchive
// +build !libarchive

package extractors

import (
	"github.com/jonathongardner/forklift/extractors/diskfs"
	"github.com/jonathongardner/forklift/extractors/gzip"
	"github.com/jonathongardner/forklift/extractors/tar"
	// "github.com/jonathongardner/forklift/extractors/gzip"
	// log "github.com/sirupsen/logrus"
)

func init() {
	toAdd := []extractor{
		//-----------------diskfs-----------------
		{diskfs.Iso9660, diskfs.ExtractArchive},
		//-----------------diskfs-----------------

		//-----------------gzip-----------------
		{"application/gzip", gzip.Extract},
		//-----------------gzip-----------------

		//-----------------tar-----------------
		{"application/x-tar", tar.Extract},
		{tar.TarGz, tar.ExtractGzip},
		//-----------------tar-----------------

	}
	for _, t := range toAdd {
		addExtractor(t)
	}
}
