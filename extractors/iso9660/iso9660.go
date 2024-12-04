package iso9660

import (
	"fmt"

	"github.com/jonathongardner/virtualfs"

	"github.com/diskfs/go-diskfs"
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/extractors/helpers"
)

const iso9660 = "application/octet-stream;iso=true"

func init() {
	isoDetector := helpers.MatchSigMultiOffsetFunc([]byte{0x43, 0x44, 0x30, 0x30, 0x31}, []int{0x8001, 0x8801, 0x9001})
	mimetype.Extend(isoDetector, iso9660, ".iso")
}
func Add(add func(string, helpers.ExtratFunc)) {
	add(iso9660, ExtractArchive)
	// add(iso9660, libarchive.ExtractArchive)
}

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	disk, err := diskfs.Open(virtualFS.LocalPath(""))
	if err != nil {
		return fmt.Errorf("diskfs error opening entry %v", err)
	}

	table, err := disk.GetPartitionTable()
	if err != nil {
		return fmt.Errorf("diskfs error opening partition table %v", err)
	}
	// Not sure if this could happen so erroring if it does, we can change if needed
	// pcount := len(table.GetPartitions())
	// if pcount != 1 {
	// 	return fmt.Errorf("unsupported multiple partitions %v", pcount)
	// }

	for index := range table.GetPartitions() {
		fs, err := disk.GetFilesystem(0) // assuming it is the whole disk, so partition = 0
		if err != nil {
			return fmt.Errorf("diskfs error opening filesystem (%v) %v", index, err)
		}
		// pvirtualFS, err := virtualFS.NewFsChild("parition-" + fmt.Sprint(index))

		err = helpers.DiskFsWalk(fs, virtualFS)
		if err != nil {
			return err
		}
	}

	return nil
}
