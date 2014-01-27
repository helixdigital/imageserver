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
	"io"
	"os"
	"strings"
	"testing"
)

var croptests = []struct {
	width          int
	height         int
	format         int
	cropTo         image.Rectangle
	expectedwidth  int
	expectedheight int
}{
	{1000, 2000, Jpg, image.Rect(0, 0, 500, 500), 500, 500},
	{1000, 2000, Gif, image.Rect(0, 0, 500, 500), 500, 500},
	{2000, 4000, Png, image.Rect(200, 200, 700, 700), 500, 500},
	{2000, 4000, Jpg, image.Rect(200, 200, 700, 700), 500, 500},
	{2000, 4000, Gif, image.Rect(100, 100, 700, 700), 600, 600},
}

func TestCropImageLoop(t *testing.T) {
	for i, tt := range croptests {
		grayimg := getGrayReader(tt.width, tt.height, tt.format)
		croprequest := CropRequest{
			grayimg,
			tt.cropTo,
			tt.format,
		}
		ireader, err := CropImage(croprequest)
		if err != nil {
			t.Errorf("%d, Crop threw error: %s", i, err)
		}
		if ireader == nil {
			t.Error("%d. Expected a reader but was nil")
		}

		config, _, err := image.DecodeConfig(ireader)
		if err != nil {
			t.Error("%d. Decode threw error: %s", i, err)
		}
		if config.Height != tt.expectedheight || config.Width != tt.expectedwidth {
			t.Errorf("%d. Resize(%d, %d) => (%d,%d), wanted (%d,%d)", i, tt.width, tt.height, config.Width, config.Height, tt.expectedwidth, tt.expectedheight)
		}
	}
}

func TestCropImage(t *testing.T) {

	// rewrite this so that we don't have files in the test suite
	inputfile, err := os.Open("/tmp/input.png")
	if err != nil {
		t.Error("Test could not open file input.png")
	}

	outputfile, err := os.Create("/tmp/cropped.png")
	if err != nil {
		t.Error("Test could not open file output.png")
	}

	croprequest := CropRequest{
		inputfile,
		image.Rect(100, 600, 500, 800),
		Png,
	}

	ireader, err := CropImage(croprequest)
	if err != nil {
		t.Error("Found error", err)
	}
	if ireader == nil {
		t.Error("Expected a reader")
	}
	_, err = io.Copy(outputfile, ireader)
	if err != nil {
		t.Error("Error when writing output file", err)
	}
}

func TestSplitAfterN(t *testing.T) {
	input := "aa.bb.cc"
	index := strings.LastIndex(input, ".")
	if input[index:] != ".cc" {
		t.Error("wanted output '.cc' but was", input[index:])
	}
}
