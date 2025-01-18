package helpers

import (
	"fmt"
	"io"
	"io/fs"
	"time"

	"github.com/jonathongardner/virtualfs"
)

func IsSymLink(mode fs.FileMode) bool {
	return (mode & fs.ModeSymlink) == fs.ModeSymlink
}

func ExtDir(virtualFS *virtualfs.Fs, name string, mode fs.FileMode, mtime time.Time) error {
	err := virtualFS.MkdirP(name, mode, mtime)
	if err != nil {
		return fmt.Errorf("couldn't extract directory %v (%v) - %v", name, virtualFS.ErrorId(), err)

	}
	return nil
}
func ExtRegular(virtualFS *virtualfs.Fs, name string, mode fs.FileMode, mtime time.Time, r io.Reader) error {
	f, err := virtualFS.Create(name, mode, mtime)
	if err != nil {
		return fmt.Errorf("couldn't create file %v (%v) - %v", name, virtualFS.ErrorId(), err)
	}
	_, err = io.Copy(f, r)
	err2 := f.Close()
	if err != nil {
		return fmt.Errorf("couldn't copy file %v %v (%v) - %v", name, mode, virtualFS.ErrorId(), err)
	}
	if err2 != nil {
		return fmt.Errorf("couldn't close file %v (%v) - %v", name, virtualFS.ErrorId(), err2)
	}
	return nil
}
func ExtSymlink(virtualFS *virtualfs.Fs, symlink, name string, mode fs.FileMode, mtime time.Time) error {
	// TODO: move | mode symlink to virtual fs
	err := virtualFS.Symlink(symlink, name, mode|fs.ModeSymlink, mtime)
	if err != nil {
		return fmt.Errorf("couldn't create symlink %v (%v) - %v", name, virtualFS.ErrorId(), err)
	}
	return nil
}
func ExtUnsuported(name string, mode fs.FileMode) error {
	return fmt.Errorf("unsupported mode %v for %v", mode, name)
}

func ExtractCompression(virtualFS *virtualfs.Fs, cf func(io.Reader) (io.Reader, error)) error {
	rf, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer rf.Close()

	fi, err := virtualFS.Stat("/")
	if err != nil {
		return fmt.Errorf("couldn't get root file stats %v - %v", virtualFS.ErrorId(), err)
	}

	f, err := virtualFS.CreateChild(fi.Mode(), fi.ModTime())
	if err != nil {
		return fmt.Errorf("couldn't create extracted file %v - %v", virtualFS.ErrorId(), err)
	}

	r, err := cf(rf)
	if err != nil {
		return fmt.Errorf("couldn't open compression file %v - %v", virtualFS.ErrorId(), err)
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
