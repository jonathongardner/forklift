package squashfs

import (
	"fmt"

	"github.com/CalebQ42/squashfs"
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"

	"github.com/jonathongardner/virtualfs"
)

const SquashFS = "application/x-squashfs-image"

func init() {
	squashFSDetector := helpers.MatchSigFunc([]byte{0x68, 0x73, 0x71, 0x73}, 0)
	mimetype.Extend(squashFSDetector, SquashFS, ".squashfs")
}

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	file, err := virtualFS.Open("/")
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fs, err := squashfs.NewReader(file)
	if err != nil {
		return fmt.Errorf("error reading squashfs: %v", err)
	}

	err = helpers.SquashFsWalk(fs, virtualFS)
	if err != nil {
		return err
	}

	return nil
}
