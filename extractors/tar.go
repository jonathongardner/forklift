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

func tarExtract(entry *fs.Entry) ([]*fs.Entry, error) {
	toReturn := make([]*fs.Entry, 0)
	file, err := entry.OpenTmp()
	if err != nil {
		return toReturn, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return toReturn, err
		}

		mode := os.FileMode(header.Mode)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := entry.ExtractedDirectory(header.Name, mode); err != nil {
				return toReturn, err
			}
		case tar.TypeReg:
			newEntry, err := entry.ExtractedFile(header.Name, mode, tarReader)
			if err != nil {
				return toReturn, err
			}

			toReturn = append(toReturn, newEntry)
		case tar.TypeSymlink:
			newEntry, err := entry.ExtractedSymlink(header.Name, mode, header.Linkname)
			if err != nil {
				return toReturn, err
			}

			toReturn = append(toReturn, newEntry)
		case tar.TypeLink:
			return toReturn, fmt.Errorf("Hardlink %v %v", header.Typeflag, header.Name)
		default:
			return toReturn, fmt.Errorf("Unknown header type %v %v", header.Typeflag, header.Name)
		}
	}
	return toReturn, nil
}

func init() {
	addExtractor(TarMtype, tarExtract)
}
