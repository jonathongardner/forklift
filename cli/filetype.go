package cli

import (
	"fmt"
	"os"

	"github.com/jonathongardner/virtualfs/filetype"

	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var filetypeCommand = &cli.Command{
	Name:      "filetype",
	Usage:     "get filetype for file",
	ArgsUsage: "[file]",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		path := c.Args().Get(0)

		reader := os.Stdin
		if path != "" {
			fileToCopy, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("Couldn't open path (%v) - %v", path, err)
			}
			defer fileToCopy.Close()
			reader = fileToCopy
		}

		ftype, err := filetype.NewFiletypeFromReader(reader)
		if err != nil {
			return fmt.Errorf("Couldn't get filetype path (%v) - %v", path, err)
		}

		fmt.Printf("File: %v, extension: %v, mimetype: %v\n", path, ftype.Extension, ftype.Mimetype)
		return nil
	},
}
