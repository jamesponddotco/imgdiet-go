package app

import (
	"fmt"
	"os"
	"path/filepath"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/urfave/cli/v2"
)

const (
	// ErrFileExists is the error returned when the file already exists.
	ErrFileExists xerrors.Error = "file already exists"

	// ErrNotEnoughArguments is the error returned when there are not enough
	// arguments.
	ErrNotEnoughArguments xerrors.Error = "not enough arguments; expected input and output"
)

// OptimizeAction is the action for the optimize command.
func OptimizeAction(c *cli.Context) error {
	if c.NArg() < 2 {
		if err := cli.ShowAppHelp(c); err != nil {
			return fmt.Errorf("%w", err)
		}

		return ErrNotEnoughArguments
	}

	var (
		input  = c.Args().Get(0)
		output = c.Args().Get(1)
		dir    = filepath.Dir(output)
		opts   = &imgdiet.Options{
			Quality:            c.Uint("quality"),
			Compression:        c.Uint("compression"),
			QuantTable:         3,
			OptimizeCoding:     true,
			Interlaced:         c.Bool("interlace"),
			StripMetadata:      c.Bool("strip"),
			OptimizeICCProfile: c.Bool("optimize-icc-profile"),
			TrellisQuant:       true,
			OvershootDeringing: true,
			OptimizeScans:      true,
		}
	)

	imgdiet.Start(nil)
	defer imgdiet.Stop()

	data, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer data.Close()

	image, err := imgdiet.Open(data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer image.Close()

	optimizedImage, err := image.Optimize(opts)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err = os.Stat(output); !os.IsNotExist(err) && !c.Bool("overwrite") {
		return fmt.Errorf("%w: %s", ErrFileExists, output)
	}

	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()

	if _, err = file.Write(optimizedImage); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
