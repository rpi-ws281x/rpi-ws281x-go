[![GoDoc](https://godoc.org/github.com/rpi-ws281x/rpi-ws281x-go?status.svg)](http://godoc.org/github.com/rpi-ws281x/rpi-ws281x-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rpi-ws281x/rpi-ws281x-go)](https://goreportcard.com/report/github.com/rpi-ws281x/rpi-ws281x-go)
[![Actions](https://github.com/rpi-ws281x/rpi-ws281x-go/workflows/CI/badge.svg)](https://github.com/rpi-ws281x/rpi-ws281x-go/actions)
[![license](https://img.shields.io/github/license/rpi-ws281x/rpi-ws281x-go.svg)](https://github.com/rpi-ws281x/rpi-ws281x-go)

# rpi-ws281x-go

## Summary

Go (golang) binding for the rpi_ws281x userspace Raspberry Pi library for controlling WS281X LEDs by Jeremy Garff ([https://github.com/jgarff/rpi_ws281x](https://github.com/jgarff/rpi_ws281x)). The goal for this library is to offer all the features of the C library and to make is as efficiently as possible.

## Installing

This module is a wrapper around the [rpi_ws281x](https://github.com/jgarff/rpi_ws281x) C library and you need to have this C library installed on your machine before installing this module.

### Compiling directly on the Raspberry Pi

**This is not the recommended way**, but if you want to compile everything on the Raspbery Pi itself, start by building
the pi_ws281x C library according to the [documentation](https://github.com/jgarff/rpi_ws281x#build),
copy the `*.a` files to `/usr/local/lib` and the `*.h` files to `/usr/local/include`.

Then you can compile the go code as usual.

### Cross compiling

The recommended way for building software for embedded systems is to use cross compilation. Cross compilation is very
easy in go... unless you use cgo. And because this module is a wrapper around a C library, we have to use cgo and the
cross compilation is not so easy.

The solution proposed here uses a docker container to cross-compile the go code and should work on GNU/Linux, macos
and Windows.

First check that you have a recent version of docker desktop. Run the following command:

```
docker buildx  ls
```

You should see something like this :

```
NAME/NODE DRIVER/ENDPOINT STATUS  PLATFORMS
default * docker
  default default         running linux/amd64, linux/arm64, ..., linux/arm/v7, linux/arm/v6
```

If you see `linux/arm/v6` and `linux/arm/v7` you can cross-compile for arm 32 bits. If you see `linux/arm64` you
can compile code for arm 64 bits.

Now you need to build a docker image with the required toolchain and libraries:

```
docker buildx build --platform linux/arm/v7 --tag ws2811-builder --file docker/app-builder/Dockerfile .
```

You can replace `linux/arm/v7` by `linux/arm64` if you want to build for arm64.

You can now use this Docker image to build your app. For example, you can build the "swiss" example using this command:

```
> cd examples/swiss
> APP="swiss"
> docker run --rm -v "$PWD":/usr/src/$APP --platform linux/arm/v7 \
  -w /usr/src/$APP ws2811-builder:latest go build -o "$APP-armv7" -v
```

On GNU/Linux or macos, you can check the built binary with the `file` command:

```
> file swiss-armv7

swiss-armv7: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV),
  dynamically linked, interpreter /lib/ld-linux-armhf.so.3, for GNU/Linux 3.2.0,
  Go BuildID=..., BuildID[sha1]=..., not stripped```
```

As you can see, the resulting binary is an executable file for the ARM processor.

## Using the module

In order to use this module, you have to understand the options of the underlying C library. Read [documentation of the C library](https://github.com/jgarff/rpi_ws281x) for more information.

The mapping of these options from go to C should be obvious. The [documentation of this module](https://godoc.org/github.com/rpi-ws281x/rpi-ws281x-go) and particularly the section about the [channel options](https://godoc.org/github.com/rpi-ws281x/rpi-ws281x-go#ChannelOption) provide further information.

## Testing

This library is tested using the following hardware setup:

<p align="center">
  <img src="https://i.imgur.com/jodJKUp.png" width="600" title="Hardware setup">
</p>

In this circuit, the 4050 is a driver that converts the 3.3V of the Raspberry Pi to the 5V needed by the ws2811 chip. The LED matrix is connected by an external power supply that provides the required current.

Here is the result of the "Swiss" example:

<p align="center">
  <img src="https://i.imgur.com/pgdvBY0.jpg" width="600" title="Swiss Example">
</p>

## Special Thanks

* Thank you to [Jeremy Garff](https://github.com/jgarff) for writing and maintaining the C library.
* Thank you to all contributors (alphabetically): 
  - [Alexandr Pavlyuk](https://github.com/pav5000)
  - [Allen Flickinger](https://github.com/FuzzyStatic) 
  - [Ben Watkins](https://github.com/OutdatedVersion)
  - [Chris C.](https://github.com/TwinProduction)
  - [Elie Grenon](https://github.com/DrunkenPoney)
  - [Herman](https://github.com/hermanbanken)
  - [Ivaylo Stoyanov](https://github.com/ivkos)
  - [Stephen Onnen](https://github.com/onnenon)

## Projects using this module

* [Rainbow and Random demo](https://github.com/FuzzyStatic/rpi-ws281x-examples-go)