/*
Copyright 2014 Helix Digital

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"image"
	"testing"
)

var imagetests = []struct {
	width          int
	height         int
	format         int
	resizetowidth  uint
	resizetoheight uint
	expectedwidth  int
	expectedheight int
}{
	{1000, 2000, Jpg, 500, 0, 500, 1000},
	{2000, 1000, Jpg, 500, 0, 500, 250},
	{2000, 1000, Jpg, 500, 0, 500, 250},
	{1000, 2000, Gif, 0, 500, 250, 500},
	{1000, 2000, Gif, 500, 500, 500, 500},
	{2000, 1000, Png, 0, 500, 1000, 500},
	{2000, 1000, Png, 0, 0, 2000, 1000},
}

func TestResizeImages(t *testing.T) {
	for i, tt := range imagetests {
		grayimg := getGrayReader(tt.width, tt.height, tt.format)
		resizerequest := ResizeRequest{
			Reader: grayimg,
			Width:  tt.resizetowidth,
			Height: tt.resizetoheight,
			Format: tt.format,
		}
		resized, err := ResizeImage(resizerequest)
		if err != nil {
			t.Errorf("%d. Resize threw error: %s", i, err)
		}
		config, _, err := image.DecodeConfig(resized)
		if err != nil {
			t.Error("%d. Decode threw error: %s", i, err)
		}
		if config.Height != tt.expectedheight || config.Width != tt.expectedwidth {
			t.Errorf("%d. Resize(%s, %s) => (%d,%d), wanted (%d,%d)", i, tt.width, tt.height, config.Width, config.Height, tt.expectedwidth, tt.expectedheight)
		}
	}
}
