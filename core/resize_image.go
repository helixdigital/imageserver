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

	"github.com/nfnt/resize"
)

// ResizeRequest is the simple data structure that describes the input to the
// ResizeImage() function
type ResizeRequest struct {
	// the input image data
	Reader io.Reader
	// The width in pixels to resize the image to.
	// Leave as 0 and set Height to keep aspect ratio
	Width uint
	// The width in pixels to resize the image to.
	// Leave as 0 and set Width to keep aspect ratio
	Height uint
	// Jpg, Gif or Png
	Format int
}

// ResizeImage returns an io.Reader storing the resized image data
func ResizeImage(req ResizeRequest) (io.Reader, error) {
	img, _, err := image.Decode(req.Reader)
	if err != nil {
		return nil, err
	}

	newimg := resize.Resize(req.Width, req.Height, img, resize.Lanczos3)

	return encodeToBuffer(newimg, req.Format), nil
}
