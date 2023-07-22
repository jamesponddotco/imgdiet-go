// Package app is the main package for the application.
package app

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/imgdiet-go/cmd/imgdiet/internal/meta"
	"github.com/urfave/cli/v2"
)

// Run is the entry point for the application.
func Run() int {
	app := cli.NewApp()
	app.Name = meta.Name
	app.Version = meta.Version
	app.Usage = meta.Description
	app.HideHelpCommand = true

	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:    "quality",
			Aliases: []string{"q"},
			Usage:   "set the quality of the output image",
			Value:   60,
		},
		&cli.UintFlag{
			Name:    "compression",
			Aliases: []string{"c"},
			Usage:   "set the compression level of the output image",
			Value:   9,
		},
		&cli.BoolFlag{
			Name:    "interlace",
			Aliases: []string{"i"},
			Usage:   "whether to interlace the output image",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "strip",
			Aliases: []string{"s"},
			Usage:   "whether to strip metadata from the output image",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "optimize-icc-profile",
			Aliases: []string{"p"},
			Usage:   "whether to optimize the ICC profile of the output image",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "overwrite",
			Aliases: []string{"w"},
			Usage:   "whether to overwrite the already existing output image",
			Value:   false,
		},
	}

	app.Action = OptimizeAction

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return 1
	}

	return 0
}
