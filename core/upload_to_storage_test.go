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
	"strings"
	"testing"

	"github.com/helixdigital/imageserver/plugin/upload"
)

func TestUploadToStorage(t *testing.T) {

	req := UploadRequest{
		strings.NewReader("Test input"),
		"text/plain",
		"teststring.txt",
	}

	mock := upload.NewMock()
	InjectUploader(mock)
	Upload(req)
	if !mock.WasCalled {
		t.Error("mock was not called")
	}
}
