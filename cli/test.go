package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var testCommand = &cli.Command{
	Name:  "test",
	Usage: "test",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		fmt.Println("test")

		return nil
	},
}
