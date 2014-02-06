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

package entities

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

type Format int

const (
	Jpg = iota
	Gif
	Png
)

type Image struct {
	Img    image.Image
	Format Format
}

// creates an io.Reader of image from entities.Image
func (self Image) Reader() io.Reader {
	output := new(bytes.Buffer)
	switch self.Format {
	case Jpg:
		jpeg.Encode(output, self.Img, &jpeg.Options{85})
	case Gif:
		gif.Encode(output, self.Img, &gif.Options{256, nil, nil})
	case Png:
		png.Encode(output, self.Img)
	}
	return output
}

// CropTo returns a copy of this image that has been cropped
// to the given dimensions
func (self Image) CropTo(bounds image.Rectangle) Image {
	dst := image.NewRGBA(bounds.Sub(bounds.Min))
	r := image.Rectangle{dst.Rect.Min, dst.Rect.Min.Add(bounds.Size())}
	draw.Draw(dst, r, self.Img, bounds.Min, draw.Src)
	return Image{Img: dst, Format: self.Format}
}

func (self Image) ResizeTo(w uint, h uint) Image {
	resized := resize.Resize(w, h, self.Img, resize.Lanczos3)
	return Image{Img: resized, Format: self.Format}
}

// NewImage taking a reader and if it correctly decodes as one of
// Jpg, Gif or Png will return an entity.Image struct
func NewImage(rdr io.Reader, format Format) (Image, error) {
	src, _, err := image.Decode(rdr)
	if err != nil {
		return Image{}, err
	}
	return Image{Img: src, Format: format}, nil
}
