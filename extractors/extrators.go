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
