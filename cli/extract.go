package cli

import (
	// "fmt"

	"github.com/jonathongardner/forklift/fin"
	"github.com/jonathongardner/forklift/fs"
	"github.com/jonathongardner/forklift/box"
	"github.com/jonathongardner/forklift/routines"

	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
)

var extractCommand =  &cli.Command{
	Name:      "extract",
	Usage:     "extract files",
	ArgsUsage: "[file]",
	Flags: []cli.Flag {
		&cli.StringFlag{
			Name:    "manifest",
			Aliases: []string{"m"},
			Usage:   "manifest output (jsonl format)",
			EnvVars: []string{"FORKLIFT_MANIFEST"},
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output folder",
			Value:   fs.Path,
			EnvVars: []string{"FORKLIFT_OUTPUT"},
		},
	},
	Action:  func(c *cli.Context) error {
		file := c.Args().Get(0)
		err := fs.SetupDir(c.String("output"), file)
		if err != nil {
			return err
		}
		defer fs.CleanupDir()

		log.Infof("Starting to unload %v to %v...", file, fs.Path)

		routineController := routines.NewController()
		err = fin.Setup(c.String("manifest"), routineController)
		if err != nil {
			return err
		}

		bar, err := box.NewBarcode(file)
		if err != nil {
			return err
		}

		routineController.Go(bar)

		return routineController.IsFinished()
	},
}
