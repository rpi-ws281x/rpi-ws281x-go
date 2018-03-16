[![GoDoc](https://godoc.org/github.com/rpi-ws281x/rpi-ws281x-go?status.svg)](http://godoc.org/github.com/rpi-ws281x/rpi-ws281x-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rpi-ws281x/rpi-ws281x-go)](https://goreportcard.com/report/github.com/rpi-ws281x/rpi-ws281x-go)
[![CircleCI](https://circleci.com/gh/rpi-ws281x/rpi-ws281x-go.svg?style=shield)](https://circleci.com/gh/rpi-ws281x/rpi-ws281x-go)
[![license](https://img.shields.io/github/license/rpi-ws281x/rpi-ws281x-go.svg)](https://github.com/rpi-ws281x/rpi-ws281x-go)

# rpi-ws281x-go

## Summary

Go (golang) binding for the rpi_ws281x userspace Raspberry Pi library for controlling WS281X LEDs by Jeremy Garff ([https://github.com/jgarff/rpi_ws281x](https://github.com/jgarff/rpi_ws281x)). The goal for this library is to offer all the features of the C library and to make is as efficiently as possible.

## Testing

This library is tested using the following hardware setup:

<p align="center">
  <img src="https://i.imgur.com/jodJKUp.png" width="600" title="Hardware setup">
</p>

In this circuit, the 4050 is a driver that convert the 3.3V of the Raspberry Pi to the 5V needed by the ws2811 chip. The led matrix is connected by an external power supply that provides the required current.

Here is the the result of the "Swiss" example:

<p align="center">
  <img src="https://i.imgur.com/pgdvBY0.jpg" width="600" title="Swiss Example">
</p>