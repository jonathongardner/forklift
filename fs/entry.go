package fs

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/jonathongardner/forklift/filetype"
)

var ErrAlreadyProcessed = fmt.Errorf("node already extracted")

type entry struct {
	Size      int64             `json:"size"`
	Type      filetype.Filetype `json:"type"`
	Md5       string            `json:"md5"`
	Sha1      string            `json:"sha1"`
	Sha256    string            `json:"sha256"`
	Sha512    string            `json:"sha512"`
	Entropy   float64           `json:"entropy"`
	Processed *atomic.Bool      `json:"extracted"`
	// SymlinkPath string            `json:"symlinkPath"`
	// Archive     bool              `json:"archive"`
}

func newEntry() *entry {
	extracted := &atomic.Bool{}
	extracted.Store(false)
	return &entry{Processed: extracted}
}

func newDirEntry() *entry {
	extracted := &atomic.Bool{}
	extracted.Store(false)
	return &entry{Type: filetype.Dir, Processed: extracted}
}

// Return old value, if old valud is true then it was already extracted
// might should return an error for that?
func (e *entry) process() error {
	if e.Processed.Swap(true) {
		return ErrAlreadyProcessed
	}
	return nil
}

type FileInfo struct {
	name      string
	size      int64
	mode      os.FileMode
	modTime   time.Time
	Md5       string
	Sha1      string
	Sha256    string
	Sha512    string
	Entropy   float64
	Type      filetype.Filetype
	Processed bool
}

func NewFileInfo(n *node) FileInfo {
	return FileInfo{
		name:      n.name,
		size:      n.ref.entry.Size,
		mode:      n.mode,
		modTime:   n.modTime,
		Md5:       n.ref.entry.Md5,
		Sha1:      n.ref.entry.Sha1,
		Sha256:    n.ref.entry.Sha256,
		Sha512:    n.ref.entry.Sha512,
		Entropy:   n.ref.entry.Entropy,
		Type:      n.ref.entry.Type,
		Processed: n.ref.entry.Processed.Load(),
	}
}

func (mfi FileInfo) Name() string {
	return mfi.name
}
func (mfi FileInfo) Size() int64 {
	return mfi.size
}
func (mfi FileInfo) Mode() os.FileMode {
	return mfi.mode
}
func (mfi FileInfo) ModTime() time.Time {
	return mfi.modTime
}
func (mfi FileInfo) IsDir() bool {
	return mfi.mode.IsDir()
}
func (mfi FileInfo) Sys() any {
	return nil
}
