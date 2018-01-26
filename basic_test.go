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

package ws2811

import (
	"flag"
	"fmt"
	"os/user"
	"runtime"
	"testing"
)

var gpioPin = flag.Int("gpio-pin", 18, "GPIO pin")
var width = flag.Int("width", 32, "LED matrix width")
var height = flag.Int("height", 8, "LED matrix height")
var brightness = flag.Int("brightness", 64, "Brightness (0-255)")

const (
	pixelColor = 255 << 16 // red
)

func TestSnake(t *testing.T) {
	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOARCH == "arm" && user.Uid != "0" {
		fmt.Println("This test requires root privilege")
		fmt.Println("Please try \"sudo go test -v\"")
		t.SkipNow()
	}

	size := *width * *height
	opt := DefaultOptions
	opt.Channels[0].Brightness = *brightness
	opt.Channels[0].LedCount = size
	opt.Channels[0].GpioPin = *gpioPin

	ws, err := MakeWS2811(&opt)
	if err != nil {
		t.Fatal(err)
	}

	err = ws.Init()
	if err != nil {
		t.Fatal(err)
	}

	bitmap := make([]uint32, size)
	for i := 0; i < size; i++ {
		if i > 0 {
			bitmap[i-1] = 0
		}
		bitmap[i] = pixelColor
		ws.SetBitmap(bitmap)
		ws.Render()
		ws.Wait()
	}

	ws.Clear()
	ws.Render()
	ws.Wait()
	ws.Fini()
}

func init() {
	flag.Parse()
}
