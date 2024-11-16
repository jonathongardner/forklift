package archiver

import (
	"context"
	"fmt"
	realFs "io/fs"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"

	"github.com/mholt/archiver/v4"
	// log "github.com/sirupsen/logrus"
)

func Add(add func(string, helpers.ExtratFunc)) {
	// Supported compression formats https://github.com/mholt/archiver/tree/master#supported-compression-formats

	supportedCompressions := []string{
		"application/x-brotli", // brotli (br)
		"application/x-bzip2",  // bzip2 (bz2)
		// flate (.zip)
		"application/gzip", // gzip (gz)
		xlz4,               // lz4
		// lzip (.lz)
		"application/x-snappy-framed", // snappy (sz)
		"application/x-xz",            // xz
		// zlib (.zz)
		// zstandard (.zst)
	}

	for _, t := range supportedCompressions {
		add(t, ExtractArchive)
	}

	// Supported archive formats https://github.com/mholt/archiver/tree/master#supported-archive-formats
	supportedArchives := []string{
		"application/zip",     // .zip
		"application/x-tar",   // .tar (including any compressed variants like .tar.gz)
		"application/vnd.rar", // .rar (read-only)
	}

	for _, t := range supportedArchives {
		add(t, ExtractArchive)
	}
}

func symlinkArchiver(toExtract *fs.Entry, path string, fileinfo realFs.FileInfo) (*fs.Entry, error) {
	mode := fileinfo.Mode()
	if mode&realFs.ModeSymlink == 0 {
		return nil, fmt.Errorf("Unknown file type %v", mode)
	}

	return toExtract.ExtractedSymlink(path, mode, fileinfo.(archiver.File).LinkTarget)
}

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	// log.Infof("extracting tar %v", entry.Name)
	f, err := virtualFS.Open("/")
	if err != nil {
		return err
	}

	fsys, err := archiver.FileSystem(context.Background(), f)
	if err != nil {
		return fmt.Errorf("error opening archive (%v) - %v", virtualFS.ErrorId(), err)
	}

	return nil
}

func init() {
	// no magic so cant do anything...
	// mimetype.Extend(matchSigFunc([]byte{0x55, 0xAA}, 510), BrotliMtype, ".br")
}
