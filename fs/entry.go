package fs

import (
  "fmt"
  "io"
  "os"
  "path/filepath"

  "github.com/jonathongardner/forklift/filetype"

	log "github.com/sirupsen/logrus"
  "github.com/google/uuid"
)

type Entry struct {
	ID          string            `json:"id"`
	ParentID    string            `json:"parentId"`
	Path        string            `json:"path"`
	SymlinkPath string            `json:"symlinkPath"`
	Size        int64             `json:"size"`
	Type        filetype.Filetype `json:"type"`
	Md5         string            `json:"md5"`
	Sha1        string            `json:"sha1"`
	Sha256      string            `json:"sha256"`
	Sha512      string            `json:"sha512"`
	Extracted   bool              `json:"extracted"`
	Mode        os.FileMode       `json:"mode"`
	tmpPath     string
}

func NewEntry(path string, mode os.FileMode, parentId string) (*Entry) {
  return &Entry{ID: uuid.New().String(), Path: path, Mode: mode, ParentID: parentId}
}

func (e *Entry) ExtractedDirectory(name string, mode os.FileMode) (error) {
  fullPath := FullPath(e.extractedPath(name))
  log.Debugf("Extracting Dir %v (%v)", fullPath, mode)
  return os.MkdirAll(fullPath, mode)
}

func (e *Entry) ExtractedFile(name string, mode os.FileMode, reader io.Reader) (*Entry, error) {
  path := e.extractedPath(name)
  log.Debugf("Extracting File %v (%v)", path, mode)
	newEnt := NewEntry(path, mode, e.ID)

  if err := newEnt.mkdirAll(); err != nil {
		return nil, err
	}

  err := newEnt.createAndSetEntryInfo(reader)
	if err != nil {
		return nil, err
	}

  return newEnt, nil
}

func (e *Entry) ExtractedSymlink(name string, mode os.FileMode, target string) (*Entry, error) {
  path := e.extractedPath(name)
  log.Debugf("Extracting Link %v (%v)", path, mode)
  newEnt := NewEntry(path, mode, e.ID)

  if err := newEnt.mkdirAll(); err != nil {
		return nil, err
	}
  // if target is absolute then take into account where this is being extracted
	if filepath.IsAbs(target) {
    // WEIRD: This is different then path, which doesnt include the `full`
		newEnt.SymlinkPath = FullPath(e.extractedPath(target))
	} else { // if its relative dont need to change anything
		newEnt.SymlinkPath = target
	}
	// os.Symlink(target, symlink)
	err := os.Symlink(newEnt.SymlinkPath, FullPath(path))
	if err != nil {
		return nil, err
	}
  return newEnt, nil
}


//-------------------Tmp--------------------
func (e *Entry) MoveToTmp() (error) {
	if e.tmpPath != "" {
		return fmt.Errorf("Already moved to tmp")
	}

	e.tmpPath = uuid.New().String()
  return os.Rename(FullPath(e.Path), TmpPath(e.tmpPath))
}

func (e *Entry) RemoveTmp() (error) {
	if e.tmpPath == "" {
		return fmt.Errorf("Already deleted tmp")
	}

  return os.Remove(TmpPath(e.tmpPath))
}

func (e *Entry) OpenTmp() (*os.File, error) {
  return os.Open(TmpPath(e.tmpPath))
}
//-------------------Tmp--------------------
