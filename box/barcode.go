package box

import (
	"fmt"

	"github.com/jonathongardner/forklift/extractors"
	"github.com/jonathongardner/forklift/fs"
	"github.com/jonathongardner/forklift/routines"
	log "github.com/sirupsen/logrus"
)

// Scan a pacakges barcode to figureout what in it
// and decide to open/extract or not
type Barcode struct {
	virtualFS *fs.Virtual
}

// This is to start scanning/extracting files. New delivery returns a barcode
// that we can decide how to handle
func NewDelivery(output, toExtract string) (*Barcode, error) {
	virtualFS, err := fs.NewVirtual(output, toExtract)
	if err != nil {
		return nil, err
	}
	// defaults
	return &Barcode{virtualFS: virtualFS}, nil
}

func (b *Barcode) VirtualFS() *fs.Virtual {
	return b.virtualFS
}

// Needed for routine runnable
// "Scan" the barcode and decide if need to extract
func (b *Barcode) Run(rc *routines.Controller) error {
	// This return error if root node already extracted, so if already extracted just move on
	err := b.virtualFS.RootProcess()
	if err != nil {
		return nil
	}

	typ := b.virtualFS.RootType()
	extFunc, ok := extractors.Functions[typ.Mimetype]
	log.Debugf("Extractions %v %v %v", ok, typ.Mimetype, b.virtualFS.RootErrorID())

	if ok {
		err := extFunc(b.virtualFS)
		if err != nil {
			return fmt.Errorf("couldn't extract entry (%v) - %v", b.virtualFS.RootErrorID(), err)
		}
		for _, childrenFS := range b.virtualFS.VirtualChildren() {
			rc.Go(&Barcode{virtualFS: childrenFS})
		}
	}
	return nil
}
