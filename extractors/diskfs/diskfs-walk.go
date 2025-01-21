package diskfs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/diskfs/go-diskfs"
	dskfile "github.com/diskfs/go-diskfs/backend/file"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/jonathongardner/forklift/extractors/helpers"
	"github.com/jonathongardner/virtualfs"
	log "github.com/sirupsen/logrus"
)

type linkable interface {
	ReadLink() (string, bool)
}

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	file, err := virtualFS.Open("/")
	if err != nil {
		return fmt.Errorf("error getting path %v", err)
	}

	disk, err := diskfs.OpenBackend(dskfile.New(file, true), diskfs.WithSectorSize(4096))
	if err != nil {
		return fmt.Errorf("error getting backend %v", err)
	}
	// set the partition table if we can
	disk.GetPartitionTable()

	if disk.Table == nil {
		fs, err := disk.GetFilesystem(0)
		if err != nil {
			return fmt.Errorf("error getting filesystem %v", err)
		}
		err = DiskFsWalk(fs, virtualFS)
		if err != nil {
			return fmt.Errorf("error walking filesystem %v", err)
		}
	}

	partitions := disk.Table.GetPartitions()

	for i, part := range partitions {
		log.Debugf("%v - %v %v %v", i, part.UUID(), part.GetStart(), part.GetSize())
		if part.GetSize() == 0 {
			continue
		}

		fs, err := disk.GetFilesystem(i)
		if err != nil {
			virtualFS.Warning(fmt.Errorf("error getting filesystem-%v %v", i, err))
			continue
		}
		// newVirtualFs, err := virtualFS.NewFsChild(partion.UUID())
		// if err != nil {
		// 	return fmt.Errorf("error getting new filesystem-%v %v", i, err)
		// }
		err = DiskFsWalk(fs, virtualFS)
		if err != nil {
			return fmt.Errorf("error walking filesystem-%v %v", i, err)
		}
	}
	return nil
	// return fmt.Errorf("unsupported type")
}

func DiskFsWalk(fs filesystem.FileSystem, virtualFS *virtualfs.Fs) error {
	return myDiskFsWalk(fs, "/", func(name string, info os.FileInfo) error {
		mode := info.Mode()
		mtime := info.ModTime()

		if info.IsDir() {
			return helpers.ExtDir(virtualFS, name, mode, mtime)
		}
		if mode.IsRegular() {
			r, err := fs.OpenFile(name, os.O_RDONLY)
			if err != nil {
				return fmt.Errorf("couldnt open file %v", err)
			}

			return helpers.ExtRegular(virtualFS, name, mode, mtime, r)
		}
		if helpers.IsCharacterDevice(mode) {
			// r, err := fs.OpenFile(name, os.O_RDONLY)
			// if err != nil {
			// 	return fmt.Errorf("couldnt open file %v", err)
			// }
			r := bytes.NewReader([]byte{})

			err := helpers.ExtRegular(virtualFS, name, mode, mtime, r)
			if err != nil {
				return err
			}

			newFs, err := virtualFS.FsFrom(name)
			if err != nil {
				return fmt.Errorf("error getting new filesystem for tags %v", err)
			}

			newFs.TagS("type", "character-device")
			return nil
		}
		if helpers.IsDevice(mode) {
			// r, err := fs.OpenFile(name, os.O_RDONLY)
			// if err != nil {
			// 	return fmt.Errorf("couldnt open file %v", err)
			// }
			r := bytes.NewReader([]byte{})

			err := helpers.ExtRegular(virtualFS, name, mode, mtime, r)
			if err != nil {
				return err
			}

			newFs, err := virtualFS.FsFrom(name)
			if err != nil {
				return fmt.Errorf("error getting new filesystem for tags %v", err)
			}

			newFs.TagS("type", "device")
			return nil
		}
		if helpers.IsSymLink(mode) {
			value, ok := info.(linkable)
			if !ok {
				return fmt.Errorf("Not able to cast as link")
			}
			symlink, ok := value.ReadLink()
			if !ok {
				return fmt.Errorf("Not able to get link")
			}

			return helpers.ExtSymlink(virtualFS, symlink, name, mode, mtime)
		}
		err := helpers.ExtUnsuported(name, mode)
		virtualFS.Warning(err)
		return nil
	})
}

func myDiskFsWalk(fs filesystem.FileSystem, root string, fn func(string, os.FileInfo) error) error {
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
			err = myDiskFsWalk(fs, path, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening %v", err)
	}

	readOnly := true
	diskG, err := diskfs.OpenBackend(dskfile.New(file, readOnly), diskfs.WithSectorSize(4096))
	if err != nil {
		return fmt.Errorf("error getting backend %v", err)
	}
	// set the partition table if we can
	diskG.GetPartitionTable()

	if diskG.Table == nil {
		fyst, err := diskG.GetFilesystem(0)
		if err != nil {
			return fmt.Errorf("error getting filesystem %v", err)
		}
		fmt.Println(path)
		readFilesystem(fyst)
	} else {
		for i, partion := range diskG.Table.GetPartitions() {
			fyst, err := diskG.GetFilesystem(i)
			if err != nil {
				return fmt.Errorf("error getting filesystem-%v %v", i, err)
			}
			fmt.Printf("%v (%v)", path, partion.UUID())
			readFilesystem(fyst)
		}
	}
	return nil
}
func readFilesystem(fyst filesystem.FileSystem) error {
	files, err := fyst.ReadDir("/")
	if err != nil {
		return fmt.Errorf("error getting dir %v", err)
	}
	for _, f := range files {
		fmt.Printf(" - %v\n", f.Name())
	}
	return nil
}
