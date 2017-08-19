// Copyright 2017 Jacques Supcik / HEIA-FR
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Interface to ws2811 chip (neopixel driver). Make sure that you have
// ws2811.h and pwm.h in a GCC include path (e.g. /usr/local/include) and
// libws2811.a in a GCC library path (e.g. /usr/local/lib).
// See https://github.com/jgarff/rpi_ws281x for instructions

package ws2811

// #cgo CFLAGS: -std=c99
// #cgo LDFLAGS: -lws2811
// #include "ws2811.go.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	// DefaultDmaNum is the default DMA number. Usually, this is 5 ob the Raspberry Pi
	DefaultDmaNum = 5
	// TargetFreq is the target frequency. It is usually 800kHz (800000), and an go as low as 400000
	TargetFreq = 800000
	StripRGB   = 0x100800 // StripRGB is the RGB Mode
	StripRBG   = 0x100008 // StripRBG is the RBG Mode
	StripGRB   = 0x081000 // StripGRB is the GRB Mode
	StripGBR   = 0x080010 // StripGBR is the GBR Mode
	StripBRG   = 0x001008 // StripBRG is the BRG Mode
	StripBGR   = 0x000810 // StripBGR is the BGR Mode
)

// Option is the arguments for creating an instance of WS2811.
type Option struct {
	Frequency  int
	DmaNum     int
	GpioPin    int
	LedCount   int
	Brightness int
	StripeType int
	Invert     bool
}

// WS2811 represent the ws2811 device
type WS2811 struct {
	dev *C.ws2811_t
}

// DefaultOptions defines sensible default options for MakeWS2811
var DefaultOptions = Option{
	Frequency:  800000,
	DmaNum:     5,
	GpioPin:    18,
	LedCount:   16,
	Brightness: 64,
	StripeType: StripGRB,
	Invert:     false,
}

var gamma8 = []uint32{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 5, 5, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10,
	10, 10, 11, 11, 11, 12, 12, 13, 13, 13, 14, 14, 15, 15, 16, 16,
	17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
	25, 26, 27, 27, 28, 29, 29, 30, 31, 32, 32, 33, 34, 35, 35, 36,
	37, 38, 39, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 50,
	51, 52, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68,
	69, 70, 72, 73, 74, 75, 77, 78, 79, 81, 82, 83, 85, 86, 87, 89,
	90, 92, 93, 95, 96, 98, 99, 101, 102, 104, 105, 107, 109, 110, 112, 114,
	115, 117, 119, 120, 122, 124, 126, 127, 129, 131, 133, 135, 137, 138, 140, 142,
	144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 164, 167, 169, 171, 173, 175,
	177, 180, 182, 184, 186, 189, 191, 193, 196, 198, 200, 203, 205, 208, 210, 213,
	215, 218, 220, 223, 225, 228, 231, 233, 236, 239, 241, 244, 247, 249, 252, 255}

func gammaCorrected(color uint32) uint32 {
	r := (color >> 16) & 0xff
	g := (color >> 8) & 0xff
	b := (color >> 0) & 0xff
	return gamma8[r]<<16 + gamma8[g]<<8 + gamma8[b]
}

// MakeWS2811 create an instance of WS2811.
func MakeWS2811(opt *Option) (ws2811 *WS2811, err error) {
	err = nil
	ws2811 = new(WS2811)
	ws2811.dev = (*C.ws2811_t)(C.malloc(C.sizeof_ws2811_t))
	if ws2811 == nil {
		err = errors.New("Unable to allocate memory")
		return
	}
	C.memset(unsafe.Pointer(ws2811.dev), 0, C.sizeof_ws2811_t)

	ws2811.dev.freq = C.uint32_t(opt.Frequency)
	ws2811.dev.dmanum = C.int(opt.DmaNum)

	ws2811.dev.channel[0].gpionum = C.int(opt.GpioPin)
	ws2811.dev.channel[0].count = C.int(opt.LedCount)
	ws2811.dev.channel[0].brightness = C.uint8_t(opt.Brightness)
	ws2811.dev.channel[0].strip_type = C.int(opt.StripeType)
	if opt.Invert {
		ws2811.dev.channel[0].invert = C.int(1)
	} else {
		ws2811.dev.channel[0].invert = C.int(0)
	}
	return
}

// Init initialize the device. It should be called only once before any other method.
func (ws2811 *WS2811) Init() error {
	res := int(C.ws2811_init(ws2811.dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.init.%d", res)
}

// Render sends a complete frame to the LED Matrix
func (ws2811 *WS2811) Render() error {
	res := int(C.ws2811_render(ws2811.dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.render.%d", res)
}

// Wait waits for render to finish. The time needed for render is given by:
// time = 1/frequency * 8 * 3 * LedCount + 0.05
// (8 is the color depth and 3 is the number of colors (LEDs) per pixel).
// See https://cdn-shop.adafruit.com/datasheets/WS2811.pdf for more details.
func (ws2811 *WS2811) Wait() error {
	res := int(C.ws2811_wait(ws2811.dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.wait.%d", res)
}

// Fini shuts down the device.
func (ws2811 *WS2811) Fini() {
	C.ws2811_fini(ws2811.dev)
}

// SetLed defines the color of a given pixel.
func (ws2811 *WS2811) SetLed(index int, value uint32) {
	C.ws2811_set_led(ws2811.dev, 0, C.int(index), C.uint32_t(gammaCorrected(value)))
}

// SetBitmap defines the color of a all pixels.
func (ws2811 *WS2811) SetBitmap(a []uint32) {
	t := make([]uint32, len(a))
	for i, color := range a {
		t[i] = gammaCorrected(color)
	}
	C.ws2811_set_bitmap(ws2811.dev, 0, unsafe.Pointer(&t[0]), C.int(len(t)*4))
}

// Clear sets all pixels to black.
func (ws2811 *WS2811) Clear() {
	C.ws2811_clear_channel(ws2811.dev, 0)
}
