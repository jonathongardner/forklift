//go:build libarchive
// +build libarchive

package extractors

import (
	"github.com/jonathongardner/forklift/extractors/compress"
	"github.com/jonathongardner/forklift/extractors/cpio"
	"github.com/jonathongardner/forklift/extractors/diskfs"
	"github.com/jonathongardner/forklift/extractors/tar"
	"github.com/jonathongardner/forklift/extractors/zip"

	// "github.com/jonathongardner/forklift/extractors/gzip"

	"github.com/jonathongardner/forklift/extractors/libarchive"
	// log "github.com/sirupsen/logrus"
)

func init() {
	toAdd := []extractor{
		//-----------------libarchive-----------------
		// compressions
		{compress.Gzip, libarchive.ExtractArchive},
		{compress.Bzip2, libarchive.ExtractArchive},
		{compress.Xz, libarchive.ExtractArchive},
		{compress.Lzip, libarchive.ExtractArchive},
		{compress.Lzma, libarchive.ExtractArchive},
		// archives
		{tar.Tar, libarchive.ExtractArchive},
		// {"application/x-pax", libarchive.ExtractArchive},
		{cpio.Cpio, libarchive.ExtractArchive},
		{zip.Zip, libarchive.ExtractArchive},
		// {"application/mtree", libarchive.ExtractArchive},
		// {"application/ar", libarchive.ExtractArchive},
		// {"application/raw", libarchive.ExtractArchive},
		{"application/x-xar", libarchive.ExtractArchive},
		// {"application/lha", libarchive.ExtractArchive},
		// {"application/lzh", libarchive.ExtractArchive},
		{"application/x-rar-compressed", libarchive.ExtractArchive},
		{"application/vnd.ms-cab-compressed", libarchive.ExtractArchive},
		{"application/x-7z-compressed", libarchive.ExtractArchive},
		{"application/warc", libarchive.ExtractArchive},
		{diskfs.Iso9660, libarchive.ExtractArchive},
		//-----------------libarchive-----------------
	}
	for _, t := range toAdd {
		addExtractor(t)
	}
}
