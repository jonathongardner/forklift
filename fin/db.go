package fin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonathongardner/forklift/fs"

	log "github.com/sirupsen/logrus"
)

func Save(output string, virtualFS *fs.Virtual) error {
	file, err := Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	count := 0
	err = virtualFS.Walk("/", func(path string, info fs.FileInfo) error {
		toSave := map[string]any{"path": path, "info": info}
		if info.Mode().IsRegular() {
			count++
		}
		jsonString, _ := json.Marshal(toSave)
		// encoder := json.NewEncoder(file)
		// encoder.Encode(toSave)
		_, err := file.Write(jsonString)
		if err != nil {
			return err
		}
		_, err = file.WriteString("\n")
		return err
	})
	if err != nil {
		return err
	}
	log.Infof("Saved %v files", count)

	return nil
}

func Create(output string) (*os.File, error) {
	dbFile := filepath.Join(output, "fin.db")
	file, err := os.Create(dbFile)
	if err != nil {
		return nil, fmt.Errorf("error opneing file %v - %v", dbFile, err)
	}

	return file, nil
}
