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

// This is the dummy version of the package so that is compiles smoothly
// on Linux, OSX or Windows. The build tag on the next lines will ensure
// that his file is not used when compiled on a Raspberry Pi.

// +build !arm

package ws2811

// WS2811 represent the ws2811 device.
type WS2811 struct {
}

// DefaultOptions defines sensible default options for MakeWS2811.
var DefaultOptions = Option{}

// MakeWS2811 create an instance of WS2811.
func MakeWS2811(opt *Option) (ws2811 *WS2811, err error) {
	ws2811 = new(WS2811)
	return ws2811, nil
}

// Init is a dummy method.
func (ws2811 *WS2811) Init() error {
	return nil
}

// Render is a dummy method.
func (ws2811 *WS2811) Render() error {
	return nil
}

// Wait is a dummy method.
func (ws2811 *WS2811) Wait() error {
	return nil
}

// Fini is a dummy method.
func (ws2811 *WS2811) Fini() {
}

// SetLed is a dummy method.
func (ws2811 *WS2811) SetLed(index int, value uint32) {
}

// SetBitmap is a dummy method.
func (ws2811 *WS2811) SetBitmap(a []uint32) {
}

// Clear is a dummy method.
func (ws2811 *WS2811) Clear() {
}
