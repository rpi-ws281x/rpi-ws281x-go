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

// +build !arm,!arm64

package ws2811

import (
	"errors"
)

// WS2811 represent the ws2811 device
type WS2811 struct {
	initialized bool
	leds        [][]uint32
	opt         *Option
}

// HwDetect gives information about the hardware
func HwDetect() HwDesc {
	return HwDesc{
		Type:          0,
		Version:       0,
		PeriphBase:    0,
		VideocoreBase: 0,
		Desc:          "DUMMY",
	}
}

// MakeWS2811 create an instance of WS2811.
func MakeWS2811(opt *Option) (*WS2811, error) {
	ws2811 := &WS2811{
		initialized: false,
		opt:         opt,
	}

	return ws2811, nil
}

// Init initialize the device. It should be called only once before any other method.
func (ws2811 *WS2811) Init() error {
	if ws2811.initialized {
		return errors.New("device already initialized")
	}

	ws2811.leds = make([][]uint32, RpiPwmChannels)

	for i := 0; i < RpiPwmChannels; i++ {
		var ledCount int
		if i < len(ws2811.opt.Channels) {
			ledCount = ws2811.opt.Channels[i].LedCount
		} else {
			ledCount = 0
		}

		ws2811.leds[i] = make([]uint32, ledCount)
	}

	return nil
}

// SetBrightness changes the brightness of a given channel. Value between 0 and 255
func (ws2811 *WS2811) SetBrightness(channel int, brightness int) {}

// Render sends a complete frame to the LED Matrix
func (ws2811 *WS2811) Render() error {
	return nil
}

// Wait waits for render to finish. The time needed for render is given by:
// time = 1/frequency * 8 * 3 * LedCount + 0.05
// (8 is the color depth and 3 is the number of colors (LEDs) per pixel).
// See https://cdn-shop.adafruit.com/datasheets/WS2811.pdf for more details.
func (ws2811 *WS2811) Wait() error {
	return nil
}

// Fini shuts down the device and frees memory.
func (ws2811 *WS2811) Fini() {
}
