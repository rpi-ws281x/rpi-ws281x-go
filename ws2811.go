// Copyright 2016 Jacques Supcik / HEIA-FR
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

// ws2811 is a Go wrapper for the rpi_ws281x library. This package is
// supposed to run a Raspberry Pi, but it also has a "dummy" version so
// that is compiles smoothly on Linux, OSX or Windows.

package ws2811

const (
	// DefaultDmaNum is the default DMA number. Usually, this is 5 ob the Raspberry Pi
	DefaultDmaNum = 5
	// TargetFreq is the target frequency. It is usually 800kHz (800000), and
	// an go as low as 400000
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
}
