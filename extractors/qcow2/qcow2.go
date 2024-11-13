package qcow2

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	// log "github.com/sirupsen/logrus"
)

const qcow2 = "application/octet-stream;qcow2=true"

func init() {
	qcow2Detector := helpers.MatchSigFunc([]byte{0x51, 0x46, 0x49, 0xfb}, 0)
	mimetype.Extend(qcow2Detector, qcow2, ".iso")
}
func Add(add func(string, helpers.ExtratFunc)) {
	// add(qcow2, ExtractArchive)
}
func ExtractArchive(virtualFS *virtualfs.Fs) error {
	return nil
}
