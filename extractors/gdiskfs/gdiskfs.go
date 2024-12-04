package gdiskfs

import (
	"fmt"

	"github.com/diskfs/go-diskfs"
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
)

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
