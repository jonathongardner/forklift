package extractors

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonathongardner/forklift/fs"
	// log "github.com/sirupsen/logrus"
)

type extratFunc func(*fs.Entry) ([]*fs.Entry, error)

var Functions = make(map[string]extratFunc)
var Types []string

// extracts to folder
func addExtractor(mtype string, ext extratFunc) {
	Functions[mtype] = ext
	Types = append(Types, mtype)
}

func mapEntriesToFiles(entry *fs.Entry, entries map[string]*fs.Entry) ([]*fs.Entry, error) {
	toReturn := make([]*fs.Entry, 0)
	root := entry.FullPath()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if root == path {
			return nil
		}

		ent, ok := entries[path]
		if !ok {
			return fmt.Errorf("entry not found for %s", path)
		}
		if info.IsDir() {
			ent.Processed = true
			ent.Extracted = true
		}
		ent.UpdateParent(entries[ent.FullPathDir()])
		toReturn = append(toReturn, ent)
		return nil
	})

	return toReturn, err
}
