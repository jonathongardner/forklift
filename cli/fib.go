package cli

import (
	// "fmt"

	"fmt"

	"github.com/jonathongardner/forklift/extractors"
	"github.com/jonathongardner/forklift/fib"
	"github.com/jonathongardner/virtualfs"
	"github.com/jonathongardner/virtualfs/filetype"

	"github.com/urfave/cli/v2"
)

var output = &cli.StringFlag{
	Name:    "output",
	Aliases: []string{"o"},
	Usage:   "output folder",
	Value:   "forklift",
	EnvVars: []string{"FORKLIFT_OUTPUT"},
}

var fibCommand = &cli.Command{
	Name:      "build",
	Usage:     "build archive from extracted folder",
	ArgsUsage: "[file]",
	Flags: []cli.Flag{
		output,
		&cli.StringSliceFlag{
			Name:    "skips",
			Usage:   "mimetypes to skip",
			EnvVars: []string{"FORKLIFT_SKIPS"},
		},
	},
	Action: func(c *cli.Context) error {
		output := c.String("output")
		fileArchive := c.Args().Get(0)

		virtualFS, err := virtualfs.NewFsFromDir(output)
		if err != nil {
			return fmt.Errorf("unable to create virtual filesystem: %v", err)
		}
		// Need to close the virtualFS so it saves the database
		defer virtualFS.Close()

		toSkip := make(map[string]bool)
		for n, _ := range extractors.Functions {
			toSkip[n] = true
		}
		// TODO: think about directories... they are so weird
		delete(toSkip, filetype.Dir.Mimetype)

		for _, n := range c.StringSlice("skips") {
			delete(toSkip, n)
		}

		return fib.Save(fileArchive, virtualFS, toSkip)
	},
}
