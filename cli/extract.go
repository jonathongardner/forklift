package cli

import (
	// "fmt"

	"fmt"

	"github.com/jonathongardner/forklift/box"
	"github.com/jonathongardner/forklift/routines"
	"github.com/jonathongardner/virtualfs"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var extractCommand = &cli.Command{
	Name:      "extract",
	Usage:     "extract files",
	ArgsUsage: "[file]",
	Flags: []cli.Flag{
		output,
	},
	Action: func(c *cli.Context) error {
		output := c.String("output")
		file := c.Args().Get(0)

		virtualFS, err := virtualfs.NewFs(output, file)
		if err != nil {
			return fmt.Errorf("unable to create virtual filesystem: %v", err)
		}
		// Need to close the virtualFS so it saves the database
		defer virtualFS.Close()

		// Create the DB to start recording info to
		// db, err := fin.CreateDB(output)
		// if err != nil {
		// 	return err
		// }

		bar, err := box.NewDelivery(virtualFS)
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

		return virtualFS.ProcessError()
	},
}
