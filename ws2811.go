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
	// RpiPwmChannels is the number of PWM channels in the Raspberry Pi
	RpiPwmChannels = 2
	// TargetFreq is the target frequency. It is usually 800kHz (800000), and an go as low as 400000
	TargetFreq = 800000
)

// StateDesc is a map from a return state to its string description.
var StateDesc = map[int]string{
	0:   "Success",
	-1:  "Generic failure",
	-2:  "Out of memory",
	-3:  "Hardware revision is not supported",
	-4:  "Memory lock failed",
	-5:  "mmap() failed",
	-6:  "Unable to map registers into userspace",
	-7:  "Unable to initialize GPIO",
	-8:  "Unable to initialize PWM",
	-9:  "Failed to create mailbox device",
	-10: "DMA error",
	-11: "Selected GPIO not possible",
	-12: "Unable to initialize PCM",
	-13: "Unable to initialize SPI",
	-14: "SPI transfer error",
}

// ChannelOption is the list of channel options
type ChannelOption struct {
	// GpioPin is the GPIO Pin with PWM alternate function, 0 if unused
	GpioPin int
	// Invert inverts output signal
	Invert bool
	// LedCount is the number of LEDs, 0 if channel is unused
	LedCount int
	// StripeType is the strip color layout -- one of WS2811StripXXX constants
	StripeType int
	// Brightness is the maximum brightness of the LEDs. Value between 0 and 255
	Brightness int
	// WShift is the white shift value
	WShift int
	// RShift is the red shift value
	RShift int
	// GShift is the green shift value
	GShift int
	// BShift is blue shift value
	BShift int
	// Gamma is the gamma correction table
	Gamma []byte
}

// Option is the list of device options
type Option struct {
	// RenderWaitTime is the time in µs before the next render can run
	RenderWaitTime int
	// Frequency is the required output frequency
	Frequency int
	// DmaNum is the number of a DMA _not_ already in use
	DmaNum int
	// Channles are channels options
	Channels []ChannelOption
}

// WS2811 represent the ws2811 device
type WS2811 struct {
	dev *C.ws2811_t
}

// DefaultOptions defines sensible default options for MakeWS2811
var DefaultOptions = Option{
	Frequency: 800000,
	DmaNum:    5,
	Channels: []ChannelOption{
		ChannelOption{
			GpioPin:    18,
			LedCount:   16,
			Brightness: 64,
			StripeType: WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
	},
}

// MakeWS2811 create an instance of WS2811.
func MakeWS2811(opt *Option) (ws2811 *WS2811, err error) {
	ws2811 = new(WS2811)
	ws2811.dev = C.ws2811_new()
	if ws2811 == nil {
		err = errors.New("Unable to allocate memory")
		return
	}
	// Reset structure
	C.memset(unsafe.Pointer(ws2811.dev), 0, C.sizeof_ws2811_t) // #nosec

	ws2811.dev.freq = C.uint32_t(opt.Frequency)
	ws2811.dev.dmanum = C.int(opt.DmaNum)

	for i, cOpt := range opt.Channels {
		c := C.ws2811_channel(ws2811.dev, C.int(i))
		_ = c // to suppress gotype warning
		c.gpionum = C.int(cOpt.GpioPin)
		c.count = C.int(cOpt.LedCount)
		c.brightness = C.uint8_t(cOpt.Brightness)
		c.strip_type = C.int(cOpt.StripeType)
		c.wshift = C.uint8_t(cOpt.WShift)
		c.rshift = C.uint8_t(cOpt.RShift)
		c.gshift = C.uint8_t(cOpt.GShift)
		c.bshift = C.uint8_t(cOpt.BShift)

		if cOpt.Invert {
			c.invert = C.int(1)
		} else {
			c.invert = C.int(0)
		}
		if cOpt.Gamma != nil {
			c.gamma = (*C.uint8_t)(&cOpt.Gamma[0])
		}
	}
	return
}

// Init initialize the device. It should be called only once before any other method.
func (ws2811 *WS2811) Init() error {
	res := int(C.ws2811_init(ws2811.dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.init: %d (%v)", res, StatusDesc(res))
}

// Render sends a complete frame to the LED Matrix
func (ws2811 *WS2811) Render() error {
	res := int(C.ws2811_render(ws2811.dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.render: %d (%v)", res, StatusDesc(res))
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
	return fmt.Errorf("Error ws2811.wait: %d (%v)", res, StatusDesc(res))
}

// Fini shuts down the device.
func (ws2811 *WS2811) Fini() {
	C.ws2811_fini(ws2811.dev)
}

// SetLed defines the color of a given pixel.
func (ws2811 *WS2811) SetLed(index int, value uint32) {
	C.ws2811_set_led(ws2811.dev, 0, C.int(index), C.uint32_t(value))
}

// SetBitmap defines the color of a all pixels.
func (ws2811 *WS2811) SetBitmap(a []uint32) {
	C.ws2811_set_bitmap(ws2811.dev, 0, unsafe.Pointer(&a[0]), C.int(len(a)*4)) // #nosec
}

// Clear sets all pixels to black.
func (ws2811 *WS2811) Clear() {
	C.ws2811_clear_channel(ws2811.dev, 0)
}

// StatusDesc returns the description of a status code
func StatusDesc(code int) string {
	desc, ok := StateDesc[code]
	if ok {
		return desc
	}
	return "unknown"
}
