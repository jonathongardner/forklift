package fs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jonathongardner/forklift/filetype"

	// log "github.com/sirupsen/logrus"
)


func (e *Entry) FullPath() string {
  return FullPath(e.Path)
}

func (e *Entry) FullPathDir() string {
  return FullPathDir(e.Path)
}
// Makes base directories
func (e *Entry) mkdirAll() (error) {
	// 0o001 allow reading in directory
  return os.MkdirAll(FullPathDir(e.Path), e.Mode) //  | 0o007
}

func (e *Entry) extractedPath(p string) string {
  return filepath.Join(e.Path, p)
}

func (e *Entry) createAndSetEntryInfo(src io.Reader) error {
  dst, err := os.OpenFile(FullPath(e.Path), os.O_CREATE|os.O_RDWR, e.Mode)
	if err != nil {
    return fmt.Errorf("Error opening - %v", err)
	}
  defer dst.Close()

	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	sha512 := sha512.New()
	ftype := filetype.NewFiletypeWriter()
	mw := io.MultiWriter(md5, sha1, sha256, sha512, ftype, dst)
	written, err := io.Copy(mw, src)
	if err != nil {
    return fmt.Errorf("Error copying - %v", err)
	}

	e.Type = ftype.String()
	e.Md5 = hex.EncodeToString(md5.Sum(nil))
	e.Sha1 = hex.EncodeToString(sha1.Sum(nil))
	e.Sha256 = hex.EncodeToString(sha256.Sum(nil))
	e.Sha512 = hex.EncodeToString(sha512.Sum(nil))
	e.Size = written

	return nil
}
