//go:build libarchive
// +build libarchive

package app

import (
	"fmt"

	"github.com/jonathongardner/libarchive"
)

Version = fmt.Sprintf("%s (libarchive: %s)", Version, libarchive.LibArchiveVersion)
