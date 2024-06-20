package box

import (
	"fmt"

	"github.com/jonathongardner/forklift/extractors"
	"github.com/jonathongardner/forklift/routines"
	"github.com/jonathongardner/virtualfs"
	log "github.com/sirupsen/logrus"
)

// Scan a pacakges barcode to figureout what in it
// and decide to open/extract or not
type Barcode struct {
	virtualFS *virtualfs.Fs
}

// This is to start scanning/extracting files. New delivery returns a barcode
// that we can decide how to handle
func NewDelivery(fs *virtualfs.Fs) (*Barcode, error) {
	return &Barcode{virtualFS: fs}, nil
}

func (b *Barcode) VirtualFS() *virtualfs.Fs {
	return b.virtualFS
}

// Needed for routine runnable
// "Scan" the barcode and decide if need to extract
func (b *Barcode) Run(rc *routines.Controller) error {
	// This return error if root node already extracted, so if already extracted just move on
	err := b.virtualFS.Process()
	if err != nil {
		return nil
	}

	fi, err := b.virtualFS.StatAt("/", 0)
	if err != nil {
		return fmt.Errorf("couldn't get stats %v - (%v)", b.virtualFS.ErrorId(), err)
	}

	typ := fi.Filetype()
	extFunc, ok := extractors.Functions[typ.Mimetype]
	log.Debugf("Extractions %v %v %v", ok, typ.Mimetype, b.virtualFS.ErrorId())

	if ok {
		b.virtualFS.Extract()
		err := extFunc(b.virtualFS)
		if err != nil {
			b.virtualFS.Error(err)
			return nil
		}
		for _, childrenFS := range b.virtualFS.FsChildren() {
			rc.Go(&Barcode{virtualFS: childrenFS})
		}
	}
	return nil
}
