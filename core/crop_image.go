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
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// This application only handles JPEG, GIF and PNG image files
const (
	Jpg = iota
	Gif
	Png
)

// CropRequest is the simple data structure that describes
// the input to the CropImage function.
type CropRequest struct {
	// The image source as an io.Reader
	Input io.Reader
	// The coordinates of the part of the Input image source
	// that the final cropped image will contain.
	Bounds image.Rectangle
	//  Jpg, Gif or Png
	Format int
}

// CropImage return an io.Reader that is the image data of the
// cropped part of the image passed in the CropRequest
func CropImage(params CropRequest) (io.Reader, error) {
	src, err := getSrcFrom(params)
	if err != nil {
		return nil, err
	}
	dst := getDstFrom(params)
	copyBoundedImage(src, dst, params.Bounds)
	return encodeToBuffer(dst, params.Format), nil
}

// converts io.Reader of image to image.Image
func getSrcFrom(params CropRequest) (image.Image, error) {
	src, _, err := image.Decode(params.Input)
	return src, err
}

// creates a new image.Image for output. Has the bounds
// as specified in the CropRequest
func getDstFrom(params CropRequest) *image.RGBA {
	return image.NewRGBA(params.Bounds.Sub(params.Bounds.Min))
}

// copies from the source to the destination
func copyBoundedImage(src image.Image, dst *image.RGBA, bounds image.Rectangle) {
	r := image.Rectangle{dst.Rect.Min, dst.Rect.Min.Add(bounds.Size())}
	draw.Draw(dst, r, src, bounds.Min, draw.Src)
}

// creates an io.Reader of image from image.Image
func encodeToBuffer(dst image.Image, format int) io.Reader {
	output := new(bytes.Buffer)
	switch format {
	case Jpg:
		jpeg.Encode(output, dst, &jpeg.Options{85})
	case Gif:
		gif.Encode(output, dst, &gif.Options{256, nil, nil})
	case Png:
		png.Encode(output, dst)
	}
	return output
}
