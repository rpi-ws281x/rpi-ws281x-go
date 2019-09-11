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

// This file contains the constants for all ws281x stripe types.

package ws2811

// 4 color R, G, B and W ordering
const (
	// SK6812StripRGBW is the RGBW Mode
	SK6812StripRGBW = 0x18100800
	// SK6812StripRBGW is the StripRBGW Mode
	SK6812StripRBGW = 0x18100008
	// SK6812StripGRBW is the StripGRBW Mode
	SK6812StripGRBW = 0x18081000
	// SK6812StrioGBRW is the StrioGBRW Mode
	SK6812StrioGBRW = 0x18080010
	// SK6812StrioBRGW is the StrioBRGW Mode
	SK6812StrioBRGW = 0x18001008
	// SK6812StripBGRW is the StripBGRW Mode
	SK6812StripBGRW = 0x18000810
	// SK6812ShiftWMask is the Shift White Mask
	SK6812ShiftWMask = 0xf0000000
)

// 3 color R, G and B ordering
const (
	// WS2811StripRGB is the RGB Mode
	WS2811StripRGB = 0x100800
	// WS2811StripRBG is the RBG Mode
	WS2811StripRBG = 0x100008
	// WS2811StripGRB is the GRB Mode
	WS2811StripGRB = 0x081000
	// WS2811StripGBR is the GBR Mode
	WS2811StripGBR = 0x080010
	// WS2811StripBRG is the BRG Mode
	WS2811StripBRG = 0x001008
	// WS2811StripBGR is the BGR Mode
	WS2811StripBGR = 0x000810
)

// Predefined fixed LED types
const (
	// WS2812Strip is the WS2812 Mode
	WS2812Strip = WS2811StripGRB
	// SK6812Strip is the SK6812 Mode
	SK6812Strip = WS2811StripGRB
	// SK6812WStrip is the SK6812W Mode
	SK6812WStrip = SK6812StripGRBW
)
