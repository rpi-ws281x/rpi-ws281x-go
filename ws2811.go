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

package ws2811

import "github.com/pkg/errors"

const (
	// DefaultDmaNum is the default DMA number.
	DefaultDmaNum = 10
	// RpiPwmChannels is the number of PWM leds in the Raspberry Pi
	RpiPwmChannels = 2
	// TargetFreq is the target frequency. It is usually 800kHz (800000), and an go as low as 400000
	TargetFreq = 800000
	// DefaultGpioPin is the default pin on the Raspberry Pi where the signal will be available. Note
	// that it is the BCM (Broadcom Pin Number) and the "Pin" 18 is actually the physical pin 12 of the
	// Raspberry Pi.
	DefaultGpioPin = 18
	// DefaultLedCount is the default number of LEDs on the stripe.
	DefaultLedCount = 16
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
//nolint: gochecknoglobals
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

// DefaultOptions defines sensible default options for MakeWS2811
//nolint: gochecknoglobals
var DefaultOptions = Option{
	Frequency: TargetFreq,
	DmaNum:    DefaultDmaNum,
	Channels: []ChannelOption{
		{
			GpioPin:    DefaultGpioPin,
			LedCount:   DefaultLedCount,
			Brightness: DefaultBrightness,
			StripeType: WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
	},
}

// Leds returns the LEDs array of a given channel
func (ws2811 *WS2811) Leds(channel int) []uint32 {
	return ws2811.leds[channel]
}

// SetLedsSync wait for the frame to finish and replace all the LEDs
func (ws2811 *WS2811) SetLedsSync(channel int, leds []uint32) error {
	if err := ws2811.Wait(); err != nil {
		return errors.WithMessage(err, "Error setting LEDs")
	}

	l := len(leds)

	if l > len(ws2811.leds[channel]) {
		return errors.New("Error: Too many LEDs")
	}

	for i := 0; i < l; i++ {
		ws2811.leds[channel][i] = leds[i]
	}

	return nil
}

// StatusDesc returns the description of a status code
func StatusDesc(code int) string {
	desc, ok := StateDesc[code]
	if ok {
		return desc
	}

	return "Unknown"
}
