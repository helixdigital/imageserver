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
	"io"
)

// UploadRequest is the simple data structure that describes the input to the
// Upload() function
type UploadRequest struct {
	// The data of the file to upload
	Reader io.Reader
	// the MimeType to tell S3 to serve the file as
	MimeType string
	// the name to store on S3 as
	UploadedName string
}

var uploader Uploader

// Uploader is the plugin that will store the image on S3 - or whatever storage provider
type Uploader interface {
	Upload([]byte, string, string) error
}

// The setter for the current uploader
func InjectUploader(upl Uploader) {
	uploader = upl
}

// Upload will store the given file on S3
func Upload(req UploadRequest) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Reader)
	if err != nil {
		return err
	}
	return uploader.Upload(buf.Bytes(), req.MimeType, req.UploadedName)
}
