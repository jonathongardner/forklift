package zip

import (
	"archive/zip"
	"fmt"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
)

const Zip = "application/zip"

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	tmpPath, err := virtualFS.Path("/")
	if err != nil {
		return err
	}

	reader, err := zip.OpenReader(tmpPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		err := extractFile(virtualFS, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(virtualFS *virtualfs.Fs, file *zip.File) error {
	fileinfo := file.FileInfo()
	mode := fileinfo.Mode()
	mtime := fileinfo.ModTime()
	if fileinfo.IsDir() {
		return helpers.ExtDir(virtualFS, file.Name, mode, mtime)
	}

	if mode.IsRegular() {
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening zip file %v (%v)", file.Name, err)
		}
		defer rc.Close()
		return helpers.ExtRegular(virtualFS, file.Name, mode, mtime, rc)
	}

	err := helpers.ExtUnsuported(file.Name, mode)
	virtualFS.Warning(err)
	return nil
}
