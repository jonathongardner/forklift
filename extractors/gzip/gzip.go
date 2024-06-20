package gzip

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	// log "github.com/sirupsen/logrus"
)

func Add(add func(string, helpers.ExtratFunc)) {
	add("application/gzip", ExtractArchive)
}

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	rf, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer rf.Close()

	r, err := gzip.NewReader(rf)
	if err != nil {
		return fmt.Errorf("couldn't open gzip file %v - %v", virtualFS.ErrorId(), err)
	}

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
