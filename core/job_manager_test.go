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
	"os"
	"testing"

	"github.com/helixdigital/imageserver/plugin/storage"
	"github.com/helixdigital/imageserver/plugin/upload"
)

func TestNewJob(t *testing.T) {
	mock := upload.NewMock()
	InjectUploader(mock)
	store := storage.NewJobStore()
	InjectJobstore(&store)

	err := MakeGrayFile(1000, 1000, "/tmp/upload.png")
	if err != nil {
		t.Error("Error in creating test file")
	}
	defer os.Remove("/tmp/upload.png")

	req := JobRequest{
		"/tmp/upload.png",
		image.Rect(200, 200, 200, 200),
		0, 150,
		"test.png",
	}

	jobid := NewJob(req)
	if jobid != 0 {
		t.Errorf("Expected jobid to be %d but was %d\n", 0, jobid)
	}

}
