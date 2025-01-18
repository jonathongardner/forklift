package cpio

import (
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"

	"kraftkit.sh/cpio"
)

const Cpio = "application/x-cpio"

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	file, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := cpio.NewReader(file)

	for {
		hdr, _, err := reader.Next()
		if err == io.EOF {
			// end of cpio archive
			break
		}
		if err != nil {
			return fmt.Errorf("error getting cpio next (%v)", err)
		}

		err = extractFile(virtualFS, hdr, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(virtualFS *virtualfs.Fs, hdr *cpio.Header, reader io.Reader) error {
	fileinfo := hdr.FileInfo()
	mode := fileinfo.Mode()
	mtime := fileinfo.ModTime()
	if fileinfo.IsDir() {
		return helpers.ExtDir(virtualFS, hdr.Name, mode, mtime)
	}

	if mode.IsRegular() {
		return helpers.ExtRegular(virtualFS, hdr.Name, mode, mtime, reader)
	}

	if helpers.IsSymLink(mode) {
		return helpers.ExtSymlink(virtualFS, hdr.Linkname, hdr.Name, mode, mtime)
	}

	err := helpers.ExtUnsuported(hdr.Name, mode)
	virtualFS.Warning(err)
	return nil
}
