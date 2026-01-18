# imgdiet

[![Go Documentation](https://godocs.io/git.sr.ht/~jamesponddotco/imgdiet-go?status.svg)](https://godocs.io/git.sr.ht/~jamesponddotco/imgdiet-go)
[![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~jamesponddotco/imgdiet-go)](https://goreportcard.com/report/git.sr.ht/~jamesponddotco/imgdiet-go)
[![Coverage Report](https://img.shields.io/badge/coverage-85%25-green)](https://git.sr.ht/~jamesponddotco/imgdiet-go/tree/trunk/item/cover.out)
[![builds.sr.ht status](https://builds.sr.ht/~jamesponddotco/imgdiet-go.svg)](https://builds.sr.ht/~jamesponddotco/imgdiet-go?)

`imgdiet` is a Go module built for optimizing images. It leverages the
power of the [`libvips`](https://github.com/libvips/libvips) library to
provide an easy-to-use, lightweight, and idiomatic way to reduce image
size without significant loss of quality.

## Prerequisites

You'll need to have `Go v1.25` and `libvips v8.18` installed on your
system to use `imgdiet`. 

## Installation

To install `imgdiet` and use it in your project, run:

```console
go get git.sr.ht/~jamesponddotco/imgdiet-go@latest
```

## Documentation

Please [see the Go reference
documentation](https://pkg.go.dev/git.sr.ht/~jamesponddotco/imgdiet-go).

## Usage

```go
package main

import (
	"context"
	"log"
	"os"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
)

func main() {
	// Start libvips with our default settings.
	imgdiet.Start(nil)
	defer imgdiet.Stop()

	// Open an image from a file as an io.Reader.
	file, err := os.Open("/path/to/image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Open the image for processing.
	img, err := imgdiet.Open(context.Background(), file)
	if err != nil {
		log.Fatal(err)
	}

	// Export the image with default optimization settings for the web, and as a
	// WebP image.
	exported, err := img.Export(imgdiet.FormatWebP, imgdiet.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	// Check to see how many bytes were saved.
	log.Printf("Saved %d bytes.", img.Size()-int64(len(exported)))
}
```

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
