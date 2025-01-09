//go:build libarchive
// +build libarchive

package extractors

import (
	"github.com/jonathongardner/forklift/extractors/diskfs"
	"github.com/jonathongardner/forklift/extractors/tar"

	// "github.com/jonathongardner/forklift/extractors/gzip"

	"github.com/jonathongardner/forklift/extractors/libarchive"
	// log "github.com/sirupsen/logrus"
)

func init() {
	toAdd := []extractor{
		//-----------------libarchive-----------------
		// compressions
		{"application/gzip", libarchive.ExtractArchive},
		{"application/x-bzip2", libarchive.ExtractArchive},
		{"application/x-xz", libarchive.ExtractArchive},
		{"application/lzip", libarchive.ExtractArchive},
		// {"application/x-lzma", libarchive.ExtractArchive},
		// archives
		{"application/x-tar", libarchive.ExtractArchive},
		{tar.TarGz, libarchive.ExtractArchive},
		// {"application/x-pax", libarchive.ExtractArchive},
		{"application/x-cpio", libarchive.ExtractArchive},
		{"application/zip", libarchive.ExtractArchive},
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
