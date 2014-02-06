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
	"fmt"
	"image"
	"io"
	"os"
	"strings"
	"time"

	"github.com/helixdigital/imageserver/entities"
)

// JobRequest is the simple data structure that describes the input to the NewJob() function
type JobRequest struct {
	// The filename of the image that will be cropped, resized and uploaded
	Local_filename string
	// the coordinates and dimensions of the part of the input image to crop to
	Crop_to image.Rectangle
	// The width in pixels to resize the cropped image to before uploading.
	// Leave as 0 and set resize_height to keep aspect ratio
	Resize_width uint
	// The height in pixels to resize the cropped image to before uploading.
	// Leave as 0 and set resize_width to keep aspect ratio
	Resize_height uint
	// The name that the cropped and resized image will be stored on S3 as.
	Uploaded_filename string
}

var jobstore entities.JobStore

// The setter for the current JobStore
func InjectJobstore(store entities.JobStore) {
	jobstore = store
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

// NewJob takes a JobRequest and starts executing
// it. It returns a jobid that can later be used
// to query the status of job's progress.
func NewJob(req JobRequest) int {
	jobid := jobstore.AssignFreeId()
	c := startOneJob(req)
	jobstore.AddJob(entities.CreateJob(jobid, c))
	go startJobWatcher(jobid)
	return jobid
}

func startJobWatcher(jobid int) {
	job, _ := jobstore.GetJob(jobid)
	for {
		select {
		case msg := <-job.Statuschan:
			switch msg.Statuscode {
			case 100:
				saveNewStatus(job, msg.Status)
			case 200:
				saveNewStatus(job, msg.Status)
				return
			case 400:
				job.Err = msg.Err
				saveNewStatus(job, msg.Status)
				return
			}
		case <-time.After(20 * time.Minute):
			job.Err = fmt.Errorf("Timed out after %s", job.Status)
			saveNewStatus(job, "Timed out")
			return
		}
	}
}

func saveNewStatus(job entities.Job, status string) {
	job.Status = status
	jobstore.Replace(job.Id, job)
}

// executes the job. Returns a channel which sends
// a msg each time the status changes
func startOneJob(req JobRequest) <-chan entities.StatusMsg {
	statuschannel := make(chan entities.StatusMsg)
	go func() {
		inputreader := readTheFile(req, statuschannel)
		defer inputreader.Close()
		original_image := getImage(req, inputreader, statuschannel)

		cropped_image := cropImage(req, original_image, statuschannel)

		resized_image := resizeImage(req, cropped_image, statuschannel)

		uploadFile(req, resized_image, statuschannel)

		statuschannel <- entities.StatusMsg{200, "Done", nil}
	}()
	return statuschannel
}

// executes the readfile part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func readTheFile(req JobRequest, statuschannel chan entities.StatusMsg) *os.File {
	statuschannel <- entities.StatusMsg{100, "Reading the file", nil}
	inputreader, err := os.Open(req.Local_filename)
	if err != nil {
		statuschannel <- entities.StatusMsg{400, "Error reading the file", err}
	}
	return inputreader
}

func getImage(req JobRequest, inputreader io.Reader, statuschannel chan entities.StatusMsg) entities.Image {
	statuschannel <- entities.StatusMsg{100, "Decoding the file", nil}
	img, err := entities.NewImage(inputreader, extension(req.Local_filename))
	if err != nil {
		statuschannel <- entities.StatusMsg{400, "Error decoding the file", err}
	}
	return img
}

// executes the cropImage part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func cropImage(req JobRequest, original_image entities.Image, statuschannel chan entities.StatusMsg) entities.Image {
	statuschannel <- entities.StatusMsg{100, "Cropping", nil}
	return original_image.CropTo(req.Crop_to)
}

// executes the resizeImage part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func resizeImage(req JobRequest, original_image entities.Image, statuschannel chan entities.StatusMsg) entities.Image {
	statuschannel <- entities.StatusMsg{100, "Resizing", nil}
	return original_image.ResizeTo(req.Resize_width, req.Resize_height)
}

// executes the uploadFile part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func uploadFile(req JobRequest, image_to_upload entities.Image, statuschannel chan entities.StatusMsg) {
	statuschannel <- entities.StatusMsg{100, "Uploading", nil}
	if err := sendToUploader(
		image_to_upload.Reader(),
		mimetype(req.Uploaded_filename),
		req.Uploaded_filename,
	); err != nil {
		statuschannel <- entities.StatusMsg{400, "Error in uploading", err}
	}
}

// Upload will store the given file on S3
func sendToUploader(rdr io.Reader, mime string, uploadedName string) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(rdr); err != nil {
		return err
	}
	return uploader.Upload(buf.Bytes(), mime, uploadedName)
}

func mimetype(filename string) string {
	ds := map[entities.Format]string{
		entities.Jpg: "image/jpeg",
		entities.Png: "image/png",
		entities.Gif: "image/gif",
	}
	return ds[extension(filename)]
}

// JobStatus returns the current status of the job with the given jobid
func JobStatus(jobid int) (string, error) {
	job, ok := jobstore.GetJob(jobid)
	if !ok {
		return fmt.Sprintf("No job found with id %d", jobid), fmt.Errorf("No job found with id %d", jobid)
	}
	return job.Status, job.Err
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
