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
	"encoding/base64"
	"image"
	_ "image/png"
	"strings"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 90
	width      = 8
	height     = 8
	ledCounts  = width * height
	maxCount   = 50
	sleepTime  = 200
)

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type invader struct {
	img     []image.Image
	current int
	ws      wsEngine
}

func (inv *invader) setup(images ...string) error {
	inv.img = make([]image.Image, len(images))
	for i, data := range images {
		var err error
		r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
		inv.img[i], _, err = image.Decode(r)
		if err != nil {
			return err
		}
	}
	inv.current = 0
	return inv.ws.Init()
}

func coordinatesToIndex(bounds image.Rectangle, x int, y int) int {
	if x%2 == 0 {
		return (x-bounds.Min.X)*height + (y - bounds.Min.Y)
	}
	return (x-bounds.Min.X)*height + (height - 1) - (y - bounds.Min.Y)
}

func rgbToColor(r uint32, g uint32, b uint32) uint32 {
	return ((r>>8)&0xff)<<16 + ((g>>8)&0xff)<<8 + ((b >> 8) & 0xff)
}

func (inv *invader) display() error {
	bounds := inv.img[inv.current].Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := inv.img[inv.current].At(x, y).RGBA()
			inv.ws.Leds(0)[coordinatesToIndex(bounds, x, y)] = rgbToColor(r, g, b)
		}
	}
	return inv.ws.Render()
}

func (inv *invader) next() {
	inv.current = (inv.current + 1) % len(inv.img)
}

func main() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	check(err)

	inv := &invader{
		ws: dev,
	}
	check(inv.setup(invader0Png, invader1Png))
	defer dev.Fini()

	for count := 0; count < maxCount; count++ {
		inv.display()
		inv.next()
		time.Sleep(sleepTime * time.Millisecond)
	}
}

const invader0Png = `
iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAACXBIWXMAABLqAAAS6gEWyM/fAAAAB3RJ
TUUH4gEfAAwq1Uz4igAAAExJREFUCNdtjLENwDAIBM9IP0T2yEBu2CsDuWEKNwyRwg1CpuGOFz8kAcCa
G3i/56jVawWrUjNrDQfW3CM8uY4kSeHZgCo146z2EZ4/za8lhsTPr+cAAAAASUVORK5CYII=
`

const invader1Png = `
iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAABmJLR0QAAAAAAAD5Q7t/AAAACXBIWXMA
ABLqAAAS6gEWyM/fAAAAB3RJTUUH4gEfCgAfIt32cwAAAFZJREFUCNdjZGVlZWBgYGBgOJPwnIGBwWSB
JITLhCyKzGBC5iDLMUEovVnCEFEI40zCc8aLqW8YsAJWVlZWVtaLqW/QGAwQV0GEIKIQLgNcFFkHAwMD
AIRcIQCYxqiMAAAAAElFTkSuQmCC
`
