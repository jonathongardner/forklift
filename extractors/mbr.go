package extractors

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/forklift/fs"
	myFs "github.com/jonathongardner/forklift/fs"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/filesystem"
	// log "github.com/sirupsen/logrus"
)

const MbrMtype = "application/octet-stream;mbr=true"

func myWalkDir(fs filesystem.FileSystem, root string, fn func(string, os.FileInfo) error) error {
	files, err := fs.ReadDir(root) // this should list everything
	if err != nil {
		return fmt.Errorf("diskfs error opening directory %v", err)
	}

	for _, file := range files {
		path := filepath.Join(root, file.Name())
		err = fn(path, file)
		if err != nil {
			return err
		}

		if file.IsDir() {
			err = myWalkDir(fs, path, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: add mbr extractors
func mbrExtract(toExtract *fs.Entry) ([]*fs.Entry, error) {
	disk, err := diskfs.Open(toExtract.TmpPath())
	if err != nil {
		return nil, fmt.Errorf("diskfs error opening entry %v", err)
	}

	table, err := disk.GetPartitionTable()
	if err != nil {
		return nil, fmt.Errorf("diskfs error opening partition table %v", err)
	}

	toReturn := make([]*myFs.Entry, 0)
	for index := range table.GetPartitions() {
		fs, err := disk.GetFilesystem(index) // assuming it is the whole disk, so partition = 0
		if err != nil {
			return nil, fmt.Errorf("diskfs error opening filesystem (%v) %v", index, err)
		}

		err = myWalkDir(fs, "/", func(p string, fileinfo os.FileInfo) error {
			mode := fileinfo.Mode()
			path := filepath.Join("parition-"+fmt.Sprint(index), p)

			switch {
			case fileinfo.IsDir():
				entry, err := toExtract.ExtractedDirectory(path, mode)
				if err != nil {
					return fmt.Errorf("diskfs error creating dir (%v) %v", index, err)
				}
				toReturn = append(toReturn, entry)
			case mode.IsRegular():
				f, err := fs.OpenFile(path, os.O_RDONLY)
				if err != nil {
					return fmt.Errorf("diskfs error opening file (%v) %v", index, err)
				}
				defer f.Close()

				entry, err := toExtract.ExtractedFile(path, mode, f)
				if err != nil {
					return fmt.Errorf("diskfs error creating file (%v) %v", index, err)
				}
				toReturn = append(toReturn, entry)
			default:
				return fmt.Errorf("Unknown file type %v", fileinfo.Mode())
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return toReturn, nil
}

func init() {
	mimetype.Extend(matchSigFunc([]byte{0x55, 0xAA}, 510), MbrMtype, ".img")
	addExtractor(MbrMtype, mbrExtract)
}
