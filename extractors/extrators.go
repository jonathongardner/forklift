package extractors

import (
	"github.com/jonathongardner/forklift/extractors/directory"
	// "github.com/jonathongardner/forklift/extractors/gzip"
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/forklift/extractors/libarchive"
	"github.com/jonathongardner/forklift/extractors/qcow2"
	// log "github.com/sirupsen/logrus"
)

var Functions = make(map[string]helpers.ExtratFunc)
var Types []string

// extracts to folder
func addExtractor(mtype string, ext helpers.ExtratFunc) {
	if _, ok := Functions[mtype]; ok {
		panic("extractor already exists")
	}
	Functions[mtype] = ext
	Types = append(Types, mtype)
}

func init() {
	directory.Add(addExtractor)
	// gzip.Add(addExtractor)
	libarchive.Add(addExtractor)
	qcow2.Add(addExtractor)
}
