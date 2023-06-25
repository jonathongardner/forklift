package extractors

import (
	// log "github.com/sirupsen/logrus"
	"github.com/gabriel-vasile/mimetype"
)

const IsoMtype = "application/octet-stream;iso=true"


func init() {
	isoDetector := matchSigMultiOffsetFunc([]byte{0x43, 0x44, 0x30, 0x30, 0x31}, []int{0x8001, 0x8801, 0x9001})
	mimetype.Extend(isoDetector, IsoMtype, ".iso")
	// addExtractor(IsoMtype, isoExtract)
}
