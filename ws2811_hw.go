// Copyright 2019 Jacques Supcik / HEIA-FR
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

// +build arm arm64

package ws2811

// #cgo CFLAGS: -std=c99 -I /usr/local/include/ws2811 -I /usr/include/ws2811
// #cgo LDFLAGS: -lws2811 -lm
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

// WS2811 represent the ws2811 device
type WS2811 struct {
	dev         *C.ws2811_t
	initialized bool
	leds        [][]uint32
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
		if length <= 0 {
			continue // spi_init does not initialize channel[1]
		}
		// convert the led C array into a golang slice:
		// https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
		// 1 << 28 is the largest pseudo-size that we can use. If we try a larger number,
		// then we get a compile error: "type [N]uint32 too large".
		ws2811.leds[i] = (*[1 << 28]uint32)(unsafe.Pointer(ledsArray))[:length:length] // nolint: gas
	}
	return nil
}

// SetBrightness changes the brightness of a given channel. Value between 0 and 255
func (ws2811 *WS2811) SetBrightness(channel int, brightness int) {
	ws2811.dev.channel[channel].brightness = C.uint8_t(brightness)
}

// SetCustomGammaFactor sets a custom Gamma correction array based on a gamma correction factor
func (ws2811 *WS2811) SetCustomGammaFactor(gammaFactor float64) {
	C.ws2811_set_custom_gamma_factor(ws2811.dev, C.double(gammaFactor))
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
