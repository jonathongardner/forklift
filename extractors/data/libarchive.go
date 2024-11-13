package libarchive

import (
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"

	"github.com/jonathongardner/virtualfs"
)

func Add(add func(string, helpers.ExtratFunc)) {
	supportedCompressions := []string{
		"application/octet-stream",
	}

	for _, t := range supportedCompressions {
		add(t, ExtractArchive)
	}
}

func ExtractArchive(virtualFS *virtualfs.Fs) error {
	// log.Infof("extracting tar %v", entry.Name)
	f, err := virtualFS.Open("/")
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

const Size = 2 ^ 16
const BufferSize = 2 ^ 24

func binWindowBytes(r io.Reader, windowFunc func([]byte) error) error {
	windowBytes := make([]byte, BufferSize)
	previousSlice := make([]byte, 0)

	for {
		// Move window
		n, err := r.Read(windowBytes)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading %v", err)
		}
		if err == io.EOF {
			break
		}
		windowBytes = windowBytes[0:n]
		windowBytes = append(previousSlice, windowBytes...)
		n = len(windowBytes)

		windowSize := Size
		if windowSize > n {
			windowSize = n
		}
		// Fill in the rest of the chunk with 0
		for i := 0; i < n-windowSize; i++ {
			windowFunc(windowBytes[i : windowSize+i])
		}
		previousSlice = windowBytes[n-windowSize-1 : n]
	}

	n := len(windowBytes)
	for i := 0; i < n; i++ {
		windowFunc(windowBytes[i:n])
	}

	return nil
}
