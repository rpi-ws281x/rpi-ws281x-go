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

// This header file is a wrapper for the rpi_ws281x library:
// https://github.com/jgarff/rpi_ws281x
 
#include <stdint.h>
#include <string.h>
#include <ws2811.h>

void ws2811_set_led(ws2811_t *ws2811, int chan, int index, uint32_t value) {
	ws2811->channel[chan].leds[index] = value;
}

void ws2811_clear_channel(ws2811_t *ws2811, int chan) {
	ws2811_channel_t *channel = &ws2811->channel[chan];
	memset(channel->leds, 0, sizeof(ws2811_led_t) * channel->count);
}

void ws2811_clear_all(ws2811_t *ws2811) {
	for (int chan = 0; chan < RPI_PWM_CHANNELS; chan++) {
		ws2811_clear_channel(ws2811, chan) {
	}
}

void ws2811_set_bitmap(ws2811_t *ws2811, int chan, void* a, int len) {
	memcpy(ws2811->channel[chan].leds, a, len);
}
