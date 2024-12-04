package extractors

import (
	"github.com/jonathongardner/forklift/extractors/directory"
	"github.com/jonathongardner/virtualfs/filetype"

	// "github.com/jonathongardner/forklift/extractors/gzip"
	"github.com/jonathongardner/forklift/extractors/gdiskfs"
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/forklift/extractors/libarchive"
	// log "github.com/sirupsen/logrus"
)

var Functions = make(map[string]helpers.ExtratFunc)
var Types []string

type extractor struct {
	mtype string
	ext   helpers.ExtratFunc
}

// extracts to folder
func addExtractor(ext extractor) {
	if _, ok := Functions[ext.mtype]; ok {
		panic("extractor already exists")
	}
	Functions[ext.mtype] = ext.ext
	Types = append(Types, ext.mtype)
}

func init() {
	toAdd := []extractor{
		//-----------------directory-----------------
		{filetype.Dir.Mimetype, directory.ExtractDir},
		//-----------------directory-----------------

		//-----------------gdiskfs-----------------
		{gdiskfs.SquashFS, gdiskfs.ExtractArchive},
		//-----------------gdiskfs-----------------

		//-----------------libarchive-----------------
		// compressions
		{"application/gzip", libarchive.ExtractArchive},
		{"application/x-bzip2", libarchive.ExtractArchive},
		{"application/x-xz", libarchive.ExtractArchive},
		{"application/lzip", libarchive.ExtractArchive},
		// {"application/x-lzma", libarchive.ExtractArchive},
		// archives
		{"application/x-tar", libarchive.ExtractArchive},
		{libarchive.TarGz, libarchive.ExtractArchive},
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
		{libarchive.Iso9660, libarchive.ExtractArchive},
		//-----------------libarchive-----------------

		//-----------------gzip-----------------
		// {"application/gzip", gzip.ExtractArchive},
		//-----------------gzip-----------------

	}
	for _, t := range toAdd {
		addExtractor(t)
	}
}
