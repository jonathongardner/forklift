package extractors

import (
	"context"
	"fmt"
	realFs "io/fs"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/fs"

	"github.com/mholt/archiver/v4"
	// log "github.com/sirupsen/logrus"
)

// Supported compression formats https://github.com/mholt/archiver/tree/master#supported-compression-formats
const BrotliMtype = "application/x-brotli"        // brotli (br)
const Bzip2Mtype = "application/x-bzip2"          // bzip2 (bz2)
const GzipMtype = "application/gzip"              // gzip (gz)
const Lz4Mtype = "application/x-lz4"              // lz4
const SnappyMtype = "application/x-snappy-framed" // snappy (sz)
const XZMtype = "application/x-xz"                // xz

// zstandard (zstd)
// Supported archive formats https://github.com/mholt/archiver/tree/master#supported-archive-formats
const ZipMtype = "application/zip"     // .zip
const TarMtype = "application/x-tar"   // .tar (including any compressed variants like .tar.gz)
const RarMtyep = "application/vnd.rar" // .rar (read-only)

func symlinkArchiver(toExtract *fs.Entry, path string, fileinfo realFs.FileInfo) (*fs.Entry, error) {
	mode := fileinfo.Mode()
	if mode&realFs.ModeSymlink == 0 {
		return nil, fmt.Errorf("Unknown file type %v", mode)
	}

	return toExtract.ExtractedSymlink(path, mode, fileinfo.(archiver.File).LinkTarget)
}

func tarZipRarEtcExtract(toExtract *fs.Entry) ([]*fs.Entry, error) {
	fsys, err := archiver.FileSystem(context.Background(), toExtract.TmpPath())
	if err != nil {
		return nil, err
	}

	return walkXtract(toExtract, fsys, symlinkArchiver)
}

func init() {
	// no magic so cant do anything...
	// mimetype.Extend(matchSigFunc([]byte{0x55, 0xAA}, 510), BrotliMtype, ".br")
	mimetype.Extend(matchSigFunc([]byte{0x04, 0x22, 0x4D, 0x18}, 0), Lz4Mtype, ".lz4")
	types := [...]string{
		BrotliMtype, Bzip2Mtype, GzipMtype, Lz4Mtype, SnappyMtype, XZMtype,
		ZipMtype, TarMtype, RarMtyep,
	}

	for _, tp := range types {
		addExtractor(tp, tarZipRarEtcExtract)
	}
}
