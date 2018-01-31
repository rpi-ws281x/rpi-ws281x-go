[![GoDoc](https://godoc.org/github.com/supcik/rpi_ws281x_go?status.svg)](http://godoc.org/github.com/supcik/rpi_ws281x_go)
[![Go Report Card](https://goreportcard.com/badge/github.com/supcik/rpi_ws281x_go)](https://goreportcard.com/report/github.com/supcik/rpi_ws281x_go)
[![CircleCI](https://circleci.com/gh/supcik/rpi_ws281x_go.svg?style=shield)](https://circleci.com/gh/supcik/rpi_ws281x_go)
[![license](https://img.shields.io/github/license/supcik/rpi_ws281x_go.svg)](https://github.com/supcik/rpi_ws281x_go)

# rpi_ws281x_go

## Summary

Go (golang) binding for the rpi_ws281x userspace Raspberry Pi library for controlling WS281X LEDs by Jeremy Garff ([https://github.com/jgarff/rpi_ws281x](https://github.com/jgarff/rpi_ws281x)). The goal for this library is to offer all the features of the C library and to make is as efficiently as possible.

## Testing

This library is tested using the following hardware setup:

<p align="center">
  <img src="https://i.imgur.com/jodJKUp.png" width="256" title="Hardware setup">
</p>

In this circuit, the 4050 is a driver that convert the 3.3V of the Raspberry Pi to the 5V needed by the ws2811 chip. The led matrix is connected by an external power supply that provides the required current.

Here is the the result of the "Swiss" example:

![Swiss Demo](https://i.imgur.com/pgdvBY0.jpg)
