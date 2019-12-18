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

package main

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 128
	width      = 8
	height     = 8
	ledCounts  = width * height
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	checkError(dev.Init())
	defer dev.Fini()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			color := uint32(0xff0000)
			if x > 2 && x < 5 && y > 0 && y < 7 {
				color = 0xffffff
			}
			if x > 0 && x < 7 && y > 2 && y < 5 {
				color = 0xffffff
			}
			dev.Leds(0)[x*height+y] = color
		}
	}
	checkError(dev.Render())

}
