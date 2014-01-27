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

package upload

import "bytes"

// MockUpload implements github.com/helixdigital/imageserver/core/Uploader
//
// It does no actual uploading. It is for testing.
type MockUpload struct {
	WasCalled     bool
	CalledData    string
	CalledMime    string
	CalledUplname string
}

// Upload mocks the Upload call and stores the parameters so that tests
// can calls that were made.
func (self *MockUpload) Upload(data []byte, mime string, uplname string) error {
	(*self).WasCalled = true
	(*self).CalledData = (bytes.NewBuffer(data)).String()
	(*self).CalledMime = mime
	(*self).CalledUplname = uplname
	return nil
}

// NewMock is the MockUpload factory
func NewMock() *MockUpload {
	return new(MockUpload)
}
