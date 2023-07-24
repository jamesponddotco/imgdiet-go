# imgdiet

[![Go Documentation](https://godocs.io/git.sr.ht/~jamesponddotco/imgdiet-go?status.svg)](https://godocs.io/git.sr.ht/~jamesponddotco/imgdiet-go)
[![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~jamesponddotco/imgdiet-go)](https://goreportcard.com/report/git.sr.ht/~jamesponddotco/imgdiet-go)
[![Coverage Report](https://img.shields.io/badge/coverage-80.5%25-green)](https://git.sr.ht/~jamesponddotco/imgdiet-go/tree/trunk/item/cover.out)
[![builds.sr.ht status](https://builds.sr.ht/~jamesponddotco/imgdiet-go.svg)](https://builds.sr.ht/~jamesponddotco/imgdiet-go?)

`imgdiet` is a Go module built for optimizing images. It leverages the
power of the [`libvips`](https://github.com/libvips/libvips) library to
provide an easy-to-use, lightweight, and idiomatic way to reduce image
size without significant loss of quality.

> **Note**: Only PNG and JPG images are supported at the moment. Support
> for more image formats is expected to be added soon. [Patches are
> welcome](https://lists.sr.ht/~jamesponddotco/imgdiet-devel).

## Prerequisites

You'll need to have `libvips` installed on your system to use `imgdiet`.
If you wish to use the command-line tool as well, you'll also need
`make` and [`scdoc`](https://git.sr.ht/~sircmpwn/scdoc) installed.

## Installation

To install `imgdiet`, run:

```sh
go get git.sr.ht/~jamesponddotco/imgdiet-go
```

You can also install the command-line application by running:

```sh
make
sudo make install
```

## Usage

### As a Go module

You can use `imgdiet` in your own Go applications like this:

```go
package main

import (
	"log"
	"os"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
)

func main() {
	// Start libvips with default settings. This is optional, but
	// recommended.
	imgdiet.Start(nil)
	defer imgdiet.Stop()

	// Open an image file as an io.Reader.
	file, err := os.Open("image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Open the image for processing.
	img, err := imgdiet.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	// Optimize the image for web use with default settings.
	image, err := img.Optimize(imgdiet.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	// Do something with the byte slice of the optimized image.
}
```

For more examples and usage details, please [check the Go reference
documentation](https://godocs.io/git.sr.ht/~jamesponddotco/imgdiet-go).

### As a CLI tool

```console
$ imgdiet --help
NAME:
   imgdiet - A CLI tool to optimize and resize images

USAGE:
   imgdiet [global options] [arguments...]

VERSION:
   0.1.0

GLOBAL OPTIONS:
   --quality value, -q value      set the quality of the output image (default: 60)
   --compression value, -c value  set the compression level of the output image (default: 9)
   --interlace, -i                whether to interlace the output image (default: false)
   --strip, -s                    whether to strip metadata from the output image (default: false)
   --optimize-icc-profile, -p     whether to optimize the ICC profile of the output image (default: false)
   --overwrite, -w                whether to overwrite the already existing output image (default: false)
   --help, -h                     show help
   --version, -v                  print the version
```

See _imgdiet(1)_ after installing for more information.

## Contributing

Anyone can help make `imgdiet` better. Send patches to the [mailing
list](https://lists.sr.ht/~jamesponddotco/imgdiet-devel) and report bugs
on the [issue tracker](https://todo.sr.ht/~jamesponddotco/imgdiet).

You must sign-off your work using `git commit --signoff`. Follow the
[Linux kernel developer's certificate of
origin](https://www.kernel.org/doc/html/latest/process/submitting-patches.html#sign-your-work-the-developer-s-certificate-of-origin)
for more details.

All contributions are made under [the MIT license](LICENSE.md).

## Resources

The following resources are available:

- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/imgdiet-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/imgdiet-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/imgdiet).

---

Released under the [MIT License](LICENSE.md).
