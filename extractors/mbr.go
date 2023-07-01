package extractors

import (
	"github.com/gabriel-vasile/mimetype"
	// log "github.com/sirupsen/logrus"
)

const MbrMtype = "application/octet-stream;mbr=true"

// TODO: add mbr extractors

func init() {
	mimetype.Extend(matchSigFunc([]byte{0x55, 0xAA}, 510), MbrMtype, ".img")
	// addExtractor(MbrMtype, mbrExtract)
}
