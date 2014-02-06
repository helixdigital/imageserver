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

package presentation

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/helixdigital/imageserver/core"
	"github.com/helixdigital/imageserver/entities"
	"github.com/helixdigital/imageserver/plugin/storage"
	"github.com/helixdigital/imageserver/plugin/upload"
)

var portnum int = 9877
var mock *upload.MockUpload

func TestAllWithSetup(t *testing.T) {
	go StartWebServer(portnum)
	time.Sleep(50 * time.Millisecond)
	mock = upload.NewMock()
	core.InjectUploader(mock)
	store := storage.NewJobStore()
	core.InjectJobstore(&store)
	core.InjectStorageReporter(&store)
	err := MakeGrayFile(1000, 1000, "/tmp/upload.gif")
	if err != nil {
		t.Error("Error in creating test file")
	}
	defer os.Remove("/tmp/upload.gif")

	testStartWebserver(t)
	testStatusOfBadJob(t)
	testRequestingNewJob(t)
	testStatusOfExistingJob(t)
	testStatsReturnsJSON(t)
}

func testStartWebserver(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", portnum))
	if err != nil {
		t.Error("Get localhost unexpectedly threw an error", err)
		return
	}
	if resp.StatusCode != 501 {
		t.Error("Get localhost expected StatusCode=501 but was", resp.StatusCode)
	}
}

func testStatusOfBadJob(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/status?jobid=10", portnum))
	if err != nil {
		t.Error("Get localhost unexpectedly threw an error", err)
		return
	}
	if resp.StatusCode != 410 {
		t.Error("Get localhost expected StatusCode=410 but was", resp.StatusCode)
	}
}

func testRequestingNewJob(t *testing.T) {
	v := getTestValues()
	v.Set("debug", "1")
	resp, err := http.PostForm(fmt.Sprintf("http://localhost:%d/request", portnum), v)
	if err != nil {
		t.Error("Get localhost unexpectedly threw an error", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Get localhost expected StatusCode=200 but was", resp.StatusCode)
	}
	body, err := getBody(resp)
	if err != nil {
		t.Error("Get response body unexpectedly threw an error", err)
		return
	}
	debug_output := `core.JobRequest{Local_filename:"/tmp/upload.gif", Crop_to:image.Rectangle{Min:image.Point{X:0, Y:0}, Max:image.Point{X:200, Y:200}}, Resize_width:0x64, Resize_height:0x0, Uploaded_filename:"uploaded.gif"}`
	if !strings.Contains(body, debug_output) {
		t.Error("Expected body to return debug info but was", body)
	}
}

func testStatusOfExistingJob(t *testing.T) {
	v := getTestValues()
	resp, _ := http.PostForm(fmt.Sprintf("http://localhost:%d/request", portnum), v)
	jobid, err := getIdFromResponse(resp)
	if err != nil {
		t.Error("Get id from response unexpectedly threw an error", err)
		return
	}

	time.Sleep(5 * time.Second)
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/status?jobid=%d", portnum, jobid))
	if err != nil {
		t.Error("Get localhost unexpectedly threw an error", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Get localhost expected StatusCode=200 but was", resp.StatusCode)
	}

	body, err := getBody(resp)
	if err != nil {
		t.Error("Get response body unexpectedly threw an error", err)
	}
	if !strings.Contains(body, "Done") {
		t.Error("Expected body to return status but was", body)
	}

	if !mock.WasCalled {
		t.Error("Did not call the mock uploader")
	}
	if mock.CalledMime != "image/gif" {
		t.Error("Did not upload an 'image/gif' file but a", mock.CalledMime)
	}
	if mock.CalledUplname != "uploaded.gif" {
		t.Error("Was supposed to upload to 'uploaded.gif' but was to", mock.CalledUplname)
	}
}

func testStatsReturnsJSON(t *testing.T) {
	resp, _ := http.Get(fmt.Sprintf("http://localhost:%d/stats", portnum))
	headers := resp.Header
	if !strings.Contains(headers.Get("Content-Type"), "application/json") {
		t.Error(
			"Content-Type should be 'application/jsan' but was:",
			headers.Get("Content-Type"),
		)
	}
}

func getIdFromResponse(resp *http.Response) (int, error) {
	stringid, err := getBody(resp)
	if err != nil {
		return -1, err
	}
	i, err := strconv.ParseInt(stringid, 10, 64)
	return int(i), err
}

func getBody(resp *http.Response) (string, error) {
	robots, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(robots), err
}

func getTestValues() url.Values {
	v := url.Values{}
	v.Set("local_filename", "/tmp/upload.gif")
	v.Set("crop_to_x", "0")
	v.Set("crop_to_y", "0")
	v.Set("crop_to_w", "200")
	v.Set("crop_to_h", "200")
	v.Set("resize_width", "100")
	v.Set("uploaded_filename", "uploaded.gif")
	return v
}

func MakeGrayFile(w int, h int, filename string) error {
	image := entities.Image{getGrayImage(w, h), extension(filename)}
	outputfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputfile.Close()
	_, err = io.Copy(outputfile, image.Reader())
	return err
}

func getGrayImage(w int, h int) *image.Gray {
	rect := image.Rect(0, 0, w, h)
	return image.NewGray(rect)
}

func extension(input string) entities.Format {
	lastdot := strings.LastIndex(input, ".")
	ext := strings.ToLower(input[lastdot:])
	switch ext {
	case ".jpeg", ".jpg":
		return entities.Jpg
	case ".gif":
		return entities.Gif
	case ".png":
		return entities.Png
	}
	return entities.Png
}
