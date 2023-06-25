package fin

import (
	"fmt"
	"os"

	"github.com/jonathongardner/forklift/routines"

	log "github.com/sirupsen/logrus"
)

type runner struct {
	entry chan []byte
	file  *os.File
}

var aRunner *runner

func Setup(path string, rc *routines.Controller) error {
	if path == "" {
		log.Debug("No manifest skipping")
		return nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("manifest already exist (%v)", path)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error opneing file (%v - %v)", path, err)
	}

	aRunner = &runner{entry: make(chan []byte), file: file}
	rc.GoBackground(aRunner)

	return nil
}

// only run one cause dont want mulitple sqlite dbs open
func (r *runner) Run(rc *routines.Controller) error {
	log.Debug("Starting fin")
	count := uint64(0)
f1:
	for {
		select {
		case e := <-r.entry:
			_, err := r.file.Write(e)
			if err != nil {
				log.Errorf("Error writing to manifest %v", err)
			}
			_, err = r.file.WriteString("\n")
			if err != nil {
				log.Errorf("Error writing to manifest %v", err)
			}

			count += 1
			if count%100 == 0 {
				log.Debugf("Processed %v", count)
			}
		case <-rc.IsDone():
			log.Debugf("Finished - Processed %v", count)
			break f1
		}
	}

	return r.file.Close()
}

func AddFile(e []byte) {
	if aRunner != nil {
		aRunner.entry <- e
	}
}
