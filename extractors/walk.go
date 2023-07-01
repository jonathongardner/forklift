package extractors

import (
	"fmt"
	realFs "io/fs"

	myFs "github.com/jonathongardner/forklift/fs"
)

type unknownTypeCallbackFunc func(*myFs.Entry, string, realFs.FileInfo) (*myFs.Entry, error)

func walkXtract(toExtract *myFs.Entry, toWalk realFs.FS, callback unknownTypeCallbackFunc) ([]*myFs.Entry, error) {
	// The filesystem should return error if something is off
	// for example github.com/mholt/archiver tar raises error if try to overwrite dir with file (or something like that)

	toReturn := make([]*myFs.Entry, 0)

	add := func(newEntry *myFs.Entry, err error) error {
		if err != nil {
			return fmt.Errorf("Error walking add %v", err)
		}

		l := len(toReturn)
		if l > 0 && toReturn[l-1].FullPath() == newEntry.FullPath() {
			toReturn[l-1] = newEntry
		} else {
			toReturn = append(toReturn, newEntry)
		}

		return nil
	}

	// walk does in lexical order
	err := realFs.WalkDir(toWalk, ".", func(path string, d realFs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Error walking %v", err)
		}
		// ignore base dir
		if path == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("Error walking info %v", err)
		}

		mode := info.Mode() //d.Type()

		switch {
		case mode.IsDir():
			return add(toExtract.ExtractedDirectory(path, mode))
		case mode.IsRegular():
			f, err := toWalk.Open(path)
			if err != nil {
				return fmt.Errorf("Error walking open %v", err)
			}

			defer f.Close()
			return add(toExtract.ExtractedFile(path, mode, f))
		default:
			return add(callback(toExtract, path, info))
		}
	})

	return toReturn, err
}
