package app

import (
	"fmt"

	"github.com/jonathongardner/libarchive"
)

var Version = fmt.Sprintf("0.0.0-beta (libarchive: %s)", libarchive.LibArchiveVersion)
