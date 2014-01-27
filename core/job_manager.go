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

		cropped_image := cropImage(req, inputreader, statuschannel)

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

// executes the cropImage part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func cropImage(req JobRequest, inputreader io.Reader, statuschannel chan entities.StatusMsg) io.Reader {
	statuschannel <- entities.StatusMsg{100, "Cropping", nil}
	croprequest := CropRequest{
		Input:  inputreader,
		Bounds: req.Crop_to,
		Format: extension(req.Local_filename),
	}
	cropped_image, err := CropImage(croprequest)
	if err != nil {
		statuschannel <- entities.StatusMsg{400, "Error in cropping", err}
	}
	return cropped_image
}

// executes the resizeImage part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func resizeImage(req JobRequest, inputreader io.Reader, statuschannel chan entities.StatusMsg) io.Reader {
	statuschannel <- entities.StatusMsg{100, "Resizing", nil}
	resizerequest := ResizeRequest{
		Reader: inputreader,
		Width:  req.Resize_width,
		Height: req.Resize_height,
		Format: extension(req.Local_filename),
	}
	resized, err := ResizeImage(resizerequest)
	if err != nil {
		statuschannel <- entities.StatusMsg{400, "Error in resizing", err}
	}
	return resized
}

// executes the uploadFile part of the job. Sends a msg on the statuschannel
// when it starts or breaks
func uploadFile(req JobRequest, inputreader io.Reader, statuschannel chan entities.StatusMsg) {
	statuschannel <- entities.StatusMsg{100, "Uploading", nil}
	uploadreq := UploadRequest{
		Reader:       inputreader,
		MimeType:     mimetype(req.Uploaded_filename),
		UploadedName: req.Uploaded_filename,
	}
	err := Upload(uploadreq)
	if err != nil {
		statuschannel <- entities.StatusMsg{400, "Error in uploading", err}
	}
}

func mimetype(filename string) string {
	ds := map[int]string{
		Jpg: "image/jpeg",
		Png: "image/png",
		Gif: "image/gif",
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

func extension(input string) int {
	lastdot := strings.LastIndex(input, ".")
	ext := strings.ToLower(input[lastdot:])
	switch ext {
	case ".jpeg", ".jpg":
		return Jpg
	case ".gif":
		return Gif
	case ".png":
		return Png
	}
	return Png
}

// for testing only.
func MakeGrayFile(w int, h int, filename string) error { /* Could probably go somewhere else */
	rdr := encodeToBuffer(getGrayImage(w, h), extension(filename))
	outputfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputfile.Close()
	_, err = io.Copy(outputfile, rdr)
	return err
}

func getGrayImage(width int, height int) *image.Gray {
	rect := image.Rect(0, 0, width, height)
	return image.NewGray(rect)
}

func getGrayReader(w int, h int, format int) io.Reader {
	return encodeToBuffer(getGrayImage(w, h), format)
}
