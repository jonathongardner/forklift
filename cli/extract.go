package cli

import (
	// "fmt"

	"fmt"
	"os"

	"github.com/jonathongardner/forklift/box"
	"github.com/jonathongardner/forklift/fin"
	"github.com/jonathongardner/forklift/routines"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var extractCommand = &cli.Command{
	Name:      "extract",
	Usage:     "extract files",
	ArgsUsage: "[file]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output folder",
			Value:   "forklift",
			EnvVars: []string{"FORKLIFT_OUTPUT"},
		},
	},
	Action: func(c *cli.Context) error {
		output := c.String("output")
		err := os.Mkdir(output, 0755)
		if err != nil {
			return fmt.Errorf("unable to create output directory: %v", err)
		}

		// Create the DB to start recording info to
		// db, err := fin.CreateDB(output)
		// if err != nil {
		// 	return err
		// }

		file := c.Args().Get(0)
		bar, err := box.NewDelivery(output, file)
		if err != nil {
			return fmt.Errorf("unable to start scanning: %v", err)
		}

		routineController := routines.NewController()
		log.Infof("Starting to unload %v to %v...", file, output)
		routineController.Go(bar)

		err = routineController.IsFinished()
		if err != nil {
			return err
		}

		return fin.Save(output, bar.VirtualFS())
	},
}
