package extractors

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/jonathongardner/forklift/fs"
	// log "github.com/sirupsen/logrus"
)

const TarMtype = "application/x-tar"

func tarExtract(toExtract *fs.Entry) ([]*fs.Entry, error) {
	file, err := toExtract.OpenTmp()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileEntryMapping := make(map[string]*fs.Entry)
	// include this one so we can set parent
	fileEntryMapping[toExtract.FullPath()] = toExtract
	tarReader := tar.NewReader(file)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		mode := os.FileMode(header.Mode)

		switch header.Typeflag {
		case tar.TypeDir:
			newEntry, err := toExtract.ExtractedDirectory(header.Name, mode)
			if err != nil {
				return nil, err
			}
			fileEntryMapping[newEntry.FullPath()] = newEntry
		case tar.TypeReg:
			newEntry, err := toExtract.ExtractedFile(header.Name, mode, tarReader)
			if err != nil {
				return nil, err
			}

			fileEntryMapping[newEntry.FullPath()] = newEntry
		case tar.TypeSymlink:
			newEntry, err := toExtract.ExtractedSymlink(header.Name, mode, header.Linkname)
			if err != nil {
				return nil, err
			}

			fileEntryMapping[newEntry.FullPath()] = newEntry
		case tar.TypeLink:
			return nil, fmt.Errorf("Hardlink %v %v", header.Typeflag, header.Name)
		default:
			return nil, fmt.Errorf("Unknown header type %v %v", header.Typeflag, header.Name)
		}
	}

	return mapEntriesToFiles(toExtract, fileEntryMapping)
}

func init() {
	addExtractor(TarMtype, tarExtract)
}
