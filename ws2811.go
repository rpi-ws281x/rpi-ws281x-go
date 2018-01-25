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
	// StripRGB is the RGB Mode
	StripRGB = 0x100800
	// StripRBG is the RBG Mode
	StripRBG = 0x100008
	// StripGRB is the GRB Mode
	StripGRB = 0x081000
	// StripGBR is the GBR Mode
	StripGBR = 0x080010
	// StripBRG is the BRG Mode
	StripBRG = 0x001008
	// StripBGR is the BGR Mode
	StripBGR = 0x000810
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
	Gamma      []byte
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
	Gamma:      gamma8,
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
	// Reset structure
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
	if opt.Gamma != nil {
		ws2811.dev.gamma = *C.uint8_t(&opt.Gamma[0])
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
