//go:build libarchive
// +build libarchive

package libarchive

import (
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"

	"github.com/jonathongardner/libarchive"
	log "github.com/sirupsen/logrus"
)

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	// log.Infof("extracting tar %v", entry.Name)
	f, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := libarchive.NewReader(f)
	if err != nil {
		return fmt.Errorf("couldn't open archive reader (%v) - %v", virtualFS.ErrorId(), err)
	}
	defer r.Close()

	for {
		header, err := r.Next()
		if err == libarchive.ErrArchiveEOF {
			break
		}
		if err != nil {
			return fmt.Errorf("couldn't get next archive value (%v) - %v", virtualFS.ErrorId(), err)
		}

		// Just a compressed file not a compressed archive
		if r.IsRaw() && header.Stat().Mode().IsRegular() {
			err = saveRaw(virtualFS, r)
		} else {
			err = saveArchive(virtualFS, header, r)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
func saveRaw(virtualFS *virtualfs.Fs, r *libarchive.Reader) error {
	fi, err := virtualFS.Stat("/")
	if err != nil {
		return fmt.Errorf("couldn't get root file stats %v - %v", virtualFS.ErrorId(), err)
	}

	f, err := virtualFS.CreateChild(fi.Mode(), fi.ModTime())
	if err != nil {
		return fmt.Errorf("couldn't create extracted file %v - %v", virtualFS.ErrorId(), err)
	}

	_, err = io.Copy(f, r)
	err2 := f.Close()
	if err != nil {
		return fmt.Errorf("couldn't copy file %v - %v", virtualFS.ErrorId(), err)
	}
	if err2 != nil {
		return fmt.Errorf("couldn't close file %v - %v", virtualFS.ErrorId(), err2)
	}

	return nil
}

func saveArchive(virtualFS *virtualfs.Fs, header libarchive.ArchiveEntry, r *libarchive.Reader) error {
	info := header.Stat()
	mode := info.Mode()
	mtime := info.ModTime()
	name := header.PathName()
	log.Debugf("extracting %v", name)

	if mode.IsDir() {
		return helpers.ExtDir(virtualFS, name, mode, mtime)
	}
	if mode.IsRegular() {
		return helpers.ExtRegular(virtualFS, name, mode, mtime, r)
	}
	if helpers.IsSymLink(mode) {
		return helpers.ExtSymlink(virtualFS, header.Symlink(), name, mode, mtime)
	}

	// err := helpers.ExtUnsuported(name, mode)
	// virtualFS.Warning(err)
	return helpers.ExtUnsuported(name, mode)
}
