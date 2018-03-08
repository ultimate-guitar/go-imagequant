## ABOUT
Go bindings for libimagequant

libimagequant is a library for lossy recompression of PNG images to reduce their filesize. It is used by the pngquant tool. This  go-imagequant project is a set of bindings for libimagequant to enable its use from the Go programming language.

This binding was written by hand. The result is somewhat more idiomatic than an automated conversion, but some  defer foo.Release() calls are required for memory management.

This project forked from [https://code.ivysaur.me/go-imagequant/]() repo.

## USAGE
Usage example is provided by a sample utility cmd/gopngquant which mimics some functionality of the upstream pngquant.

The sample utility has the following options:

```
Usage of gopngquant:
  -In string
        Input filename
  -Out string
        Output filename
  -Speed int
        Speed (1 slowest, 10 fastest) (default 3)
  -Version
```

## BUILDING
Build via go build.  
This is a CGO package and requires a C compiler installed. However, if you use go install then future invocations of go build do not require the C compiler to be present.

The imagequant.go file also declares a number of CFLAGS for GCC that allow the included libimagequant (2.8 git-a425e83) to build in an optimal way without using the upstream configure/make scripts.

## LICENSE
I am releasing this binding under the ISC license, however, libimagequant itself is released under GPLv3-or-later and/or commercial licenses. You must comply with the terms of such a license when using this binding in a Go project.

## CHANGELOG

See [CHANGELOG](CHANGELOG.md).
