package libarchive

import (
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"

	"github.com/jonathongardner/forklift/fs"

	"github.com/jonathongardner/libarchive"
	log "github.com/sirupsen/logrus"
)

func Add(add func(string, helpers.ExtratFunc)) {
	// supportedCompressions := []string{
	// 	// "application/gzip",
	// 	// "application/x-bzip2",
	// 	// "application/x-xz",
	// 	// "application/lzip",
	// 	// // "application/x-lzma"
	// }

	// for _, t := range supportedCompressions {
	// 	add(t, ExtractCompression)
	// }

	supportedArchives := []string{
		"application/x-tar",
		// "application/x-pax",
		"application/x-cpio",
		iso9660,
		"application/zip",
		// "application/mtree",
		// "application/ar",
		// "application/raw",
		"application/x-xar",
		// "application/lha",
		// "application/lzh",
		"application/x-rar-compressed",
		"application/vnd.ms-cab-compressed",
		"application/x-7z-compressed",
		"application/warc",
	}

	for _, t := range supportedArchives {
		add(t, ExtractArchive)
	}
}

func ExtractArchive(virtualFS *fs.Virtual) error {
	// log.Infof("extracting tar %v", entry.Name)
	f, err := virtualFS.RootOpen()
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := libarchive.NewReader(f)
	if err != nil {
		return fmt.Errorf("couldn't open archive reader (%v) - %v", virtualFS.RootErrorID(), err)
	}
	defer r.Close()

	for {
		header, err := r.Next()
		if err == libarchive.ErrArchiveEOF {
			break
		}
		if err != nil {
			return fmt.Errorf("couldn't get next archive value (%v) - %v", virtualFS.RootErrorID(), err)
		}

		info := header.Stat()
		mode := info.Mode()
		name := header.PathName()
		log.Debugf("extracting %v", name)

		if mode.IsDir() {
			err := virtualFS.MkdirP(name, mode)
			if err != nil {
				return fmt.Errorf("couldn't extract directory %v (%v) - %v", name, virtualFS.RootErrorID(), err)

			}
		} else if mode.IsRegular() {
			f, err := virtualFS.Create(name, mode)
			if err != nil {
				return fmt.Errorf("couldn't create file %v (%v) - %v", name, virtualFS.RootErrorID(), err)
			}
			_, err = io.Copy(f, r)
			err2 := f.Close()
			if err != nil {
				return fmt.Errorf("couldn't copy file %v (%v) - %v", name, virtualFS.RootErrorID(), err)
			}
			if err2 != nil {
				return fmt.Errorf("couldn't close file %v (%v) - %v", name, virtualFS.RootErrorID(), err2)
			}
		} else {
			return fmt.Errorf("unsupported mode %v for %v", mode, name)
		}
	}

	return nil
}
