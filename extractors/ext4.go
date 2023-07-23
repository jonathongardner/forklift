package extractors

import (
	"fmt"
	"io"
	realFs "io/fs"

	"github.com/jonathongardner/forklift/fs"

	"github.com/masahiro331/go-ext4-filesystem/ext4"
	// log "github.com/sirupsen/logrus"
)

const Ext4Mtype = "application/x-brotli" // brotli (br)

func callbackExt4(toExtract *fs.Entry, path string, fileinfo realFs.FileInfo) (*fs.Entry, error) {
	return nil, fmt.Errorf("Unknown file type %v", fileinfo.Mode())
}

type ext4Cache[K string, V any] struct {
	cache map[K]V
}

func (c *ext4Cache[K, V]) Add(k K, v V) bool {
	c.cache[k] = v
	return true
}

func (c *ext4Cache[K, V]) Get(k K) (v V, evicted bool) {
	v, evicted = c.cache[k]
	return
}

func ext4Extract(toExtract *fs.Entry) ([]*fs.Entry, error) {
	f, err := toExtract.OpenTmp()
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	cache := &ext4Cache[string, any]{cache: make(map[string]any)}

	filesystem, err := ext4.NewFS(*io.NewSectionReader(f, 0, info.Size()), cache)
	if err != nil {
		return nil, err
	}

	return walkXtract(toExtract, filesystem, callbackExt4)
}

func init() {
	// no magic so cant do anything...
	// mimetype.Extend(matchSigFunc([]byte{0x55, 0xAA}, 510), Ext4Mtype, ".ext4")
	addExtractor(Ext4Mtype, ext4Extract)
}
