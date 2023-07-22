imgdiet(1)

# NAME

imgdiet - A CLI tool to optimize and resize images

# SYNOPSIS

*imgdiet* [options...] INPUT OUTPUT

# DESCRIPTION

*imgdiet* is an easy-to-use command-line tool that offers a fast and simple
image processing and compression solution, i.e., it tries to reduce the size of
an image without significant loss of quality.

Optimization attempts are not guaranteed to succeed.

# FILES

The input file is an image file encoded in either the PNG or JPG formats.

# OPTIONS

*-q*, *--quality* n
	Set the maximum quality of the output image. n is 0 (worse) to 100 (best).
	Defaults to 60.

*-c*, *--compression* n
	Set the compression level of the output image the input image is in the PNG
	format. n is 0 (minimal effort) to 9 (maximum effort). Defaults to 9.

*-i*, *--interlace*
	Whether the image should be non-interlaced (i.e., progressive-scanned) or
	interlaced. Defaults to non-interlaced.

*-s*, *--strip*
	Whether the image should have its metadata stripped. Defaults to false.

*-p*, *--optimize-icc-profile*
	Whether the image should have its ICC profile data optimized. Defaults to
	false.

*-w*, *--overwrite*
	Whether to overwrite an already existing image. Defaults to false.

*-h*, *--help*
	Show help message and quit.

*-v*, *--version*
	Show version number and quit.

# EXAMPLES

*Example 1. Optimize file with default settings*
	The following command line optimizes image "/path/to/image/file.jpg" and
	outputs the result to "/path/to/image/optimized-file.jpg" with default
	settings.

	imgdiet '/path/to/image/file.jpg' '/path/to/image/optimized-file.jpg'

*Example 2. Optimize file with custom settings*
	The following command line optimizes image "/path/to/image/file.png" while
	stripping its metadata information, and outputs the result to
	"/path/to/image/optimized-file.png".

	imgdiet -s '/path/to/image/file.png' '/path/to/image/optimized-file.png'

# REPORTING BUGS

Report bugs via email to <~jamesponddotco/imgdiet@todo.sr.ht> or via the web
interface at https://todo.sr.ht/~jamesponddotco/imgdiet.

# AUTHORS

Maintained by James Pond <james@cipher.host>.
