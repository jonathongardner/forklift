package fs

import (
  "os"
  "path"
  "path/filepath"
)

var Path = "forklift/"
var tmpDir = ""

func SetupDir(p string, file string) error {
  Path = p
  err := os.MkdirAll(Path, os.FileMode(0777))
  if err != nil {
    return err
  }

  tmpDir, err = os.MkdirTemp(Path, filepath.Base(file))
  return err
}

func CleanupDir() error {
  return os.Remove(tmpDir)
}

func FullPath(p string) string {
  // TODO: Make sure that path doesnt move outside of this basePath
  return filepath.Join(Path, p)
}

func FullPathDir(p string) string {
  return path.Dir(FullPath(p))
}

func TmpPath(p string) string {
  if tmpDir == "" {
    panic("Tmp dir not set")
  }
  return filepath.Join(tmpDir, p)
}
