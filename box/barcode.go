package box

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonathongardner/forklift/extractors"
	"github.com/jonathongardner/forklift/fin"
	"github.com/jonathongardner/forklift/fs"
	"github.com/jonathongardner/forklift/routines"

	log "github.com/sirupsen/logrus"
)

type Barcode struct {
	entry *fs.Entry
}

func NewBarcode(pathToExtract string) (*Barcode, error) {
	// defaults
	reader := os.Stdin
	mode := fs.DIR_RWX
	path := "forklift"

	if pathToExtract != "" {
		path = filepath.Base(pathToExtract)
		fileToCopy, err := os.Open(pathToExtract)
		if err != nil {
			return nil, fmt.Errorf("Couldn't open path (%v) - %v", pathToExtract, err)
		}
		defer fileToCopy.Close()

		fileInfo, err := fileToCopy.Stat()
		if err != nil {
			return nil, fmt.Errorf("Couldn't get path info (%v) - %v", pathToExtract, err)
		}
		mode = mode | fileInfo.Mode()

		if fileInfo.IsDir() {
			return nil, fmt.Errorf("Must provide a file (not a directory)")
		}
		reader = fileToCopy
	}

	parentEnt := &fs.Entry{Path: ""}
	ent, err := parentEnt.ExtractedFile(path, mode, reader)
	if err != nil {
		return nil, fmt.Errorf("Error extracting file - %v", err)
	}
	return &Barcode{entry: ent}, nil
}

func (b *Barcode) addToManifest() {
	jsonByte, _ := json.Marshal(b.entry)
	fin.AddFile(jsonByte)
}

// Needed for routine runnable
func (b *Barcode) Run(rc *routines.Controller) error {
	defer b.addToManifest()
	e := b.entry

	if e.Type.Mimetype == "" {
		// We should set the initial type when we copy the file
		panic("Types cant be blank")
	}

	extFunc, ok := extractors.Functions[e.Type.Mimetype]

	if ok {
		err := e.MoveToTmp()
		if err != nil {
			return fmt.Errorf("Error moving files to tmp %v %v - %v", e.Path, e.Type, err)
		}

		entries, err := extFunc(e)
		if err != nil {
			return fmt.Errorf("Error extraction file %v %v - %v", e.Path, e.Type, err)
		}
		for _, entry := range entries {
			rc.Go(&Barcode{entry: entry})
		}

		err = e.RemoveTmp()
		if err != nil {
			return fmt.Errorf("Error deleting tmp file %v %v - %v", e.Path, e.Type, err)
		}

		e.Extracted = true
	} else {
		log.Debugf("Unknown type %v %v", e.Path, e.Type)
	}

	return nil
}
