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
copy the `*.a` files to `/usr/local/bin` and the `*.h` files to `/usr/local/include`.

Then you can compile the go code as usual.

### Cross compiling

The recommended way for building software for embedded systems is to use cross compilation. Cross compilation is very
easy in go... unless you use cgo. And because this module is a wrapper around a C library, we have to use cgo and the
cross compilation is not so easy.

The solution proposed here uses a docker container to cross-compile the go code and should work on GNU/Linux, macos
and Windows. The docker container uses a [balena](https://www.balena.io/) image that emulate an ARM processor using QEMU.

Start by writing the following `Dockerfile` (also available in the repository):

```
FROM balenalib/raspberrypi3-golang:latest-build AS builder
RUN [ "cross-build-start" ]
WORKDIR /tmp
RUN apt-get update -y && apt-get install -y scons
RUN git clone https://github.com/jgarff/rpi_ws281x.git && \
  cd rpi_ws281x && \
  scons
RUN [ "cross-build-end" ]

FROM balenalib/raspberrypi3-golang:latest
RUN [ "cross-build-start" ]
COPY --from=builder /tmp/rpi_ws281x/*.a /usr/local/lib/
COPY --from=builder /tmp/rpi_ws281x/*.h /usr/local/include/
RUN go get -v -u github.com/rpi-ws281x/rpi-ws281x-go
RUN [ "cross-build-end" ]
```

You might want to change the base image if you are not using the Raspberry Pi 3.

Now build the image with the command :

```
docker build --tag rpi-ws281x-go-builder .
```

The resulting image (`rpi-ws281x-go-builder`) contains the C library (in `/usr/local`) and a compiled version of the wrapper.
You can now use this image in a container to build your application. For example, if you want to build
the "swiss" example, run the following command:

```
docker run --rm -ti -v "$(pwd)"/examples/swiss:/go/src/swiss rpi-ws281x-go-builder /usr/bin/qemu-arm-static /bin/sh -c "go build -o src/swiss/swiss -v swiss"
```

On GNU/Linux or macos, you can check the built binary with the `file` command:

```
file examples/swiss/swiss

examples/swiss/swiss: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux-armhf.so.3, for GNU/Linux 3.2.0, BuildID[sha1]=7178d110f504aceb3fb184ec984e402fd2c8712e, not stripped
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
* Thank you to [Herman](https://github.com/hermanbanken) for his contribution to the documentation and for the idea of using balena for cross-compilation.