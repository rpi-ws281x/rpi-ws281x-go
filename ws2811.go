// Copyright 2018 Jacques Supcik / HEIA-FR
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
// #include <stdint.h>
// #include <stdlib.h>
// #include <string.h>
// #include <ws2811.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	// DefaultDmaNum is the default DMA number. Usually, this is 5 ob the Raspberry Pi
	DefaultDmaNum = 5
	// RpiPwmChannels is the number of PWM leds in the Raspberry Pi
	RpiPwmChannels = 2
	// TargetFreq is the target frequency. It is usually 800kHz (800000), and an go as low as 400000
	TargetFreq = 800000
	// DefaultGpioPin is the default pin on the Raspberry Pi where the signal will be available. Note
	// that it is the BCM (Broadcom Pin Number) and the "Pin" 18 is actually the physical pin 12 of the
	// Raspberry Pi.
	DefaultGpioPin    = 18
	// DefaultLedCount is the default number of LEDs on the stripe.
	DefaultLedCount   = 16
	// DefaultBrightness is the default maximum brightness of the LEDs. The brightness value can be between 0 and 255.
	// If the brightness is too low, the LEDs remain dark. If the brightness is too high, the system needs too much
	// current.
	DefaultBrightness = 64 // Safe value between 0 and 255.
)

const (
	// HwVerTypeUnknown represents unknown hardware
	HwVerTypeUnknown = 0
	// HwVerTypePi1 represents the Raspberry Pi 1
	HwVerTypePi1 = 1
	// HwVerTypePi2 represents the Raspberry Pi 2
	HwVerTypePi2 = 2
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

// HwDesc is the Hardware Description
type HwDesc struct {
	Type          uint32
	Version       uint32
	PeriphBase    uint32
	VideocoreBase uint32
	Desc          string
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
	// Channels are channel options
	Channels []ChannelOption
}

// WS2811 represent the ws2811 device
type WS2811 struct {
	dev         *C.ws2811_t
	initialized bool
	leds        [][]uint32
}

// DefaultOptions defines sensible default options for MakeWS2811
var DefaultOptions = Option{
	Frequency: TargetFreq,
	DmaNum:    DefaultDmaNum,
	Channels: []ChannelOption{
		ChannelOption{
			GpioPin:    DefaultGpioPin,
			LedCount:   DefaultLedCount,
			Brightness: DefaultBrightness,
			StripeType: WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
	},
}

// HwDetect gives information about the hardware
func HwDetect() HwDesc {
	hw := unsafe.Pointer(C.rpi_hw_detect()) // nolint: gas
	return HwDesc{
		Type:          uint32((*C.rpi_hw_t)(hw)._type),
		Version:       uint32((*C.rpi_hw_t)(hw).hwver),
		PeriphBase:    uint32((*C.rpi_hw_t)(hw).periph_base),
		VideocoreBase: uint32((*C.rpi_hw_t)(hw).videocore_base),
		Desc:          C.GoString((*C.rpi_hw_t)(hw).desc),
	}
}

// MakeWS2811 create an instance of WS2811.
func MakeWS2811(opt *Option) (ws2811 *WS2811, err error) {
	ws2811 = &WS2811{
		initialized: false,
	}
	if ws2811 == nil {
		err = errors.New("unable to allocate memory")
		return nil, err
	}
	// Allocate and reset structure
	ws2811.dev = (*C.ws2811_t)(C.malloc(C.sizeof_ws2811_t))
	C.memset(unsafe.Pointer(ws2811.dev), 0, C.sizeof_ws2811_t) // nolint: gas

	ws2811.dev.freq = C.uint32_t(opt.Frequency)
	ws2811.dev.dmanum = C.int(opt.DmaNum)

	for i, cOpt := range opt.Channels { // nolint: gotype
		ws2811.dev.channel[i].gpionum = C.int(cOpt.GpioPin)
		ws2811.dev.channel[i].count = C.int(cOpt.LedCount)
		ws2811.dev.channel[i].brightness = C.uint8_t(cOpt.Brightness)
		ws2811.dev.channel[i].strip_type = C.int(cOpt.StripeType)
		ws2811.dev.channel[i].wshift = C.uint8_t(cOpt.WShift)
		ws2811.dev.channel[i].rshift = C.uint8_t(cOpt.RShift)
		ws2811.dev.channel[i].gshift = C.uint8_t(cOpt.GShift)
		ws2811.dev.channel[i].bshift = C.uint8_t(cOpt.BShift)

		if cOpt.Invert {
			ws2811.dev.channel[i].invert = C.int(1)
		} else {
			ws2811.dev.channel[i].invert = C.int(0)
		}
		if cOpt.Gamma != nil {
			// allocate and copy gamma table. The memory will be freed by C.ws2811_fini().
			m := (*C.uint8_t)(C.malloc(C.size_t(256)))
			ws2811.dev.channel[i].gamma = m
			C.memcpy(unsafe.Pointer(m), unsafe.Pointer(&cOpt.Gamma[0]), C.size_t(256)) // nolint: gas
		}
	}
	return ws2811, err
}

// Init initialize the device. It should be called only once before any other method.
func (ws2811 *WS2811) Init() error {
	if ws2811.initialized {
		return errors.New("device already initialized")
	}
	res := int(C.ws2811_init(ws2811.dev))
	if res != 0 {
		return fmt.Errorf("error ws2811.init: %d (%v)", res, StatusDesc(res))
	}
	ws2811.initialized = true
	ws2811.leds = make([][]uint32, RpiPwmChannels)
	for i := 0; i < RpiPwmChannels; i++ {
		// var ledsArray *C.ws2811_led_t = C.ws2811_leds(ws2811.dev, C.int(i))
		ledsArray := ws2811.dev.channel[i].leds    // nolint: gotype
		length := int(ws2811.dev.channel[i].count) // nolint: gotype
		// convert the led C array into a golang slice:
		// https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
		// 1 << 28 is the largest pseudo-size that we can use. If we try a larger number,
		// then we get a compile error: "type [N]uint32 too large".
		ws2811.leds[i] = (*[1 << 28]uint32)(unsafe.Pointer(ledsArray))[:length:length] // nolint: gas
	}
	return nil
}

// Render sends a complete frame to the LED Matrix
func (ws2811 *WS2811) Render() error {
	res := int(C.ws2811_render(ws2811.dev))
	if res != 0 {
		return fmt.Errorf("error ws2811.render: %d (%v)", res, StatusDesc(res))
	}
	return nil
}

// Wait waits for render to finish. The time needed for render is given by:
// time = 1/frequency * 8 * 3 * LedCount + 0.05
// (8 is the color depth and 3 is the number of colors (LEDs) per pixel).
// See https://cdn-shop.adafruit.com/datasheets/WS2811.pdf for more details.
func (ws2811 *WS2811) Wait() error {
	res := int(C.ws2811_wait(ws2811.dev))
	if res != 0 {
		return fmt.Errorf("error ws2811.wait: %d (%v)", res, StatusDesc(res))

	}
	return nil
}

// Fini shuts down the device and frees memory.
func (ws2811 *WS2811) Fini() {
	C.ws2811_fini(ws2811.dev)
	// release the memory allocated by MakeWS2811. Note that we should not release
	// ws2811.dev.channel[i].gamma (also allocated by MakeWS2811) because C.ws2811_fini
	// already releases this data.
	C.free(unsafe.Pointer(ws2811.dev)) // nolint: gas
	ws2811.initialized = false
}

// Leds returns the LEDs array of a given channel
func (ws2811 *WS2811) Leds(channel int) []uint32 {
	return ws2811.leds[channel]
}

// StatusDesc returns the description of a status code
func StatusDesc(code int) string {
	desc, ok := StateDesc[code]
	if ok {
		return desc
	}
	return "Unknown"
}
