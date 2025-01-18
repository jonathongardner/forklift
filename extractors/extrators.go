package extractors

import (
	"github.com/jonathongardner/forklift/extractors/compress"
	"github.com/jonathongardner/forklift/extractors/cpio"
	"github.com/jonathongardner/forklift/extractors/directory"
	"github.com/jonathongardner/forklift/extractors/tar"
	"github.com/jonathongardner/forklift/extractors/zip"
	"github.com/jonathongardner/virtualfs/filetype"

	// "github.com/jonathongardner/forklift/extractors/gzip"
	"github.com/jonathongardner/forklift/extractors/diskfs"
	"github.com/jonathongardner/forklift/extractors/helpers"
	// log "github.com/sirupsen/logrus"
)

var Functions = make(map[string]helpers.ExtratFunc)

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
}

func init() {
	toAdd := []extractor{
		//-----------------directory-----------------
		{filetype.Dir.Mimetype, directory.ExtractDir},
		//-----------------directory-----------------

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
		//-----------------tar-----------------

		//-----------------cpio-----------------
		{cpio.Cpio, cpio.ExtractArchive},
		//-----------------cpio-----------------

		//-----------------zip-----------------
		{zip.Zip, zip.ExtractArchive},
		//-----------------zip-----------------

		//-----------------diskfs-----------------
		{diskfs.SquashFS, diskfs.ExtractArchive},
		{diskfs.Iso9660, diskfs.ExtractArchive},
		//-----------------diskfs-----------------
	}
	for _, t := range toAdd {
		addExtractor(t)
	}
}
