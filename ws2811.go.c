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

// This implements the wrapper for the rpi_ws281x library:
// https://github.com/jgarff/rpi_ws281x
 
#include <ws2811.go.h>

// ws2811_leds returns a reference (address) of the LEDs array of a channel.
ws2811_led_t* ws2811_leds(const ws2811_t* ws2811, int chan) {
	return ws2811->channel[chan].leds;
}

// ws2811_leds returns a reference (address) of the LEDs array of a channel.
int ws2811_leds_count(const ws2811_t* ws2811, int chan) {
	return ws2811->channel[chan].count;
}
