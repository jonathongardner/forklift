package tar

import (
	"archive/tar"
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
)

const Tar = "application/x-tar"

func Extract(virtualFS *virtualfs.Fs) error {
	f, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer f.Close()

	return extract(virtualFS, f)
}

func extract(virtualFS *virtualfs.Fs, reader io.Reader) error {
	tarReader := tar.NewReader(reader)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("could not get next header tar (%v)", err)
		}

		name := header.Name
		mode := header.FileInfo().Mode()
		mtime := header.ModTime

		switch header.Typeflag {
		case tar.TypeDir:
			err := helpers.ExtDir(virtualFS, name, mode, mtime)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := helpers.ExtRegular(virtualFS, name, mode, mtime, tarReader)
			if err != nil {
				return err
			}
		case tar.TypeSymlink:
			err := helpers.ExtSymlink(virtualFS, header.Linkname, name, mode, mtime)
			if err != nil {
				return err
			}
		default:
			err := helpers.ExtUnsuported(name, mode)
			virtualFS.Warning(err)
		}

	}
	return nil
}
