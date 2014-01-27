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

import (
	"net/http"
	"testing"
)

var RUNNING_AGAINST_PRODUCTION bool = false

// Rename back to TestS3Upload when you want to run against live S3
func TestS3Upload(t *testing.T) {
	if !RUNNING_AGAINST_PRODUCTION {
		return
	}
	s3 := NewAmazonS3Upload("", "", "")
	err := s3.Upload([]byte("one two three four five"), "plain/text", "test.txt")
	if err != nil {
		t.Error("S3 upload failed with error", err)
	}
	defer s3.Delete("test.txt")
	resp, err := http.Get("http://s3-ap-southeast-2.amazonaws.com/comet.is/test.txt")
	if err != nil {
		t.Error("Getting test file from S3 production failed with error: ", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Expected status code to be 200 but was", resp.StatusCode)
	}
	if resp.ContentLength != 23 {
		t.Error("Expected content length to be 23 but was", resp.ContentLength)
	}
}

func TestS3Delete(t *testing.T) {
	if !RUNNING_AGAINST_PRODUCTION {
		return
	}
	s3 := NewAmazonS3Upload("", "", "")
	err := s3.Upload([]byte("one two three four five"), "plain/text", "test.txt")
	if err != nil {
		t.Error("S3 upload failed with error", err)
	}
	s3.Delete("test.txt")
	if err != nil {
		t.Error("S3 delete failed with error", err)
	}
	resp, err := http.Get("http://s3-ap-southeast-2.amazonaws.com/comet.is/test.txt")
	if err != nil {
		t.Error("Getting test file from S3 production failed with error: ", err)
	}
	if resp.StatusCode != 404 {
		t.Error("Expected status code to be 404 but was", resp.StatusCode)
	}
}
