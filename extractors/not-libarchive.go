//go:build !libarchive
// +build !libarchive

package extractors

import (
	"github.com/jonathongardner/forklift/extractors/compress"
	"github.com/jonathongardner/forklift/extractors/cpio"
	"github.com/jonathongardner/forklift/extractors/diskfs"
	"github.com/jonathongardner/forklift/extractors/tar"
	"github.com/jonathongardner/forklift/extractors/zip"
	// log "github.com/sirupsen/logrus"
)

func init() {
	toAdd := []extractor{
		//-----------------compress-----------------
		{compress.Gzip, compress.ExtractGzip},
		{compress.Bzip2, compress.ExtractBzip2},
		{compress.Xz, compress.ExtractXz},
		{compress.Lzip, compress.ExtractLzip},
		{compress.Lzma, compress.ExtractLzma},
		//-----------------compress-----------------

		//-----------------tar-----------------
		{tar.Tar, tar.Extract},
		// TODO: think about other compressions... might be able to do it in compression with peek
		{tar.TarGz, tar.ExtractGzip},
		//-----------------tar-----------------

		//-----------------cpio-----------------
		{cpio.Cpio, cpio.ExtractArchive},
		//-----------------cpio-----------------

		//-----------------zip-----------------
		{zip.Zip, zip.ExtractArchive},
		//-----------------zip-----------------

		//-----------------diskfs-----------------
		{diskfs.Iso9660, diskfs.ExtractArchive},
		//-----------------diskfs-----------------

	}
	for _, t := range toAdd {
		addExtractor(t)
	}
}
