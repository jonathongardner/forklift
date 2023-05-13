package extractors

import (
	"github.com/jonathongardner/forklift/fs"
	// log "github.com/sirupsen/logrus"
)

type extratFunc func(*fs.Entry) ([]*fs.Entry, error)

var Functions = make(map[string]extratFunc)
var Types []string

// extracts to folder
func addExtractor(mtype string, ext extratFunc) {
	Functions[mtype] = ext
	Types = append(Types, mtype)
}

//---------------Sig--------------
func matchSig(raw []byte, toMatch []byte, offset int) bool {
	if (len(raw) < offset + len(toMatch)) {
		return false
	}

	for i := 0; i < len(toMatch); i++ {
		if (raw[offset + i] != toMatch[i]) {
			return false
		}
	}

	return true
}

func matchSigFunc(toMatch []byte, offset int) (func(raw []byte, limit uint32) bool) {
	return func(raw []byte, limit uint32) bool {
		return matchSig(raw, toMatch, offset)
	}
}

func matchSigMultiOffsetFunc(toMatch []byte, offsets []int) (func(raw []byte, limit uint32) bool) {
	return func(raw []byte, limit uint32) bool {
		for i := 0; i < len(offsets); i++ {
			if matchSig(raw, toMatch, offsets[i]) {
				return true
			}
		}
		return false
	}
}
//---------------Sig--------------
