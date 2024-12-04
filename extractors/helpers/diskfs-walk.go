package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/jonathongardner/virtualfs"
)

const MbrMtype = "application/octet-stream;mbr=true"

func DiskFsWalk(fs filesystem.FileSystem, virtualFS *virtualfs.Fs) error {
	return myWalk(fs, "/", func(name string, info os.FileInfo) error {
		mode := info.Mode()
		mtime := info.ModTime()

		if info.IsDir() {
			return ExtDir(virtualFS, name, mode, mtime)
		}
		if mode.IsRegular() {
			r, err := fs.OpenFile(name, os.O_RDONLY)
			if err != nil {
				return fmt.Errorf("couldnt open file %v", err)
			}

			return ExtRegular(virtualFS, name, mode, mtime, r)
		}
		err := ExtUnsuported(name, mode)
		virtualFS.Warning(err)
		return nil
	})
}

func myWalk(fs filesystem.FileSystem, root string, fn func(string, os.FileInfo) error) error {
	files, err := fs.ReadDir(root) // this should list everything
	if err != nil {
		return fmt.Errorf("diskfs error opening directory %v", err)
	}

	for _, file := range files {
		path := filepath.Join(root, file.Name())
		err = fn(path, file)
		if err != nil {
			return err
		}

		if file.IsDir() {
			err = myWalk(fs, path, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
