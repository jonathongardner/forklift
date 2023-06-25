package extractors

import (
	"github.com/gabriel-vasile/mimetype"
	// log "github.com/sirupsen/logrus"
)

const MbrMtype = "application/octet-stream;mbr=true"

func init() {
	mbrDetector := matchSigFunc([]byte{0x55, 0xAA}, 510) // 0x55AA
	mimetype.Extend(mbrDetector, MbrMtype, ".img")
	// addExtractor(MbrMtype, mbrExtract)
}
