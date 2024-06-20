package fib

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/jonathongardner/virtualfs"
)

func Save(output string, virtualFS *virtualfs.Fs, mp map[string]bool) error {
	file, err := create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	err = virtualFS.Walk("/", func(path string, info *virtualfs.FileInfo) error {
		_, ok := mp[info.Filetype().Mimetype]
		if ok {
			return nil
		}

		// Create a tar Header from the FileInfo data
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// Use full path as name (FileInfoHeader only takes the basename)
		// If we don't do this the directory strucuture would
		// not be preserved
		// https://golang.org/src/archive/tar/common.go?#L626
		header.Name = path

		// Write file header to the tar archive
		err = tw.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("Failed to write header %v %v", path, err)
		}

		if info.IsDir() {
			return nil
		}

		// Open the file which will be written into the archive
		file, err := info.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy file content to tar archive
		_, err = io.Copy(tw, file)

		return err
	})

	return err
}

func create(outputFile string) (*os.File, error) {
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("error opneing file %v - %v", outputFile, err)
	}

	return file, nil
}
