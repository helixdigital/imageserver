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

// Entities are the basic data elements with their basic operations.
package entities

import "time"

// Messages that each job returns to notify listeners of its status
type StatusMsg struct {
	Statuscode int
	Status     string
	Err        error
}

// JobStore is the plugin that provides a job API in front of the database
type JobStore interface {
	/* Should not be in entities but it is here because this is my clumsy way
	   of resolving a dependency cycle -- PB
	*/

	AddJob(Job)
	AssignFreeId() int
	GetJob(int) (Job, bool)
	Replace(int, Job)
}

// Job encapsulates the concept of performing the series of cropping,
// resizing and uploading tasks on one given file.
type Job struct {
	// the jobid given to this job
	Id int
	// the channel on which various tasks will report their status updates
	Statuschan <-chan StatusMsg
	// the human-readable description of the most-recently reported status of this job
	Status string
	// if Status starts with the substring "Error" then Err contains the binary error and `nil` otherwise
	Err error
	// when this job was first created
	Created time.Time
	// the time of the most recent change to this data structure
	Modified time.Time
}

// Returns a Job datastructure initialised with defaults plus the
// `id` and `c` parameters
func CreateJob(id int, c <-chan StatusMsg) Job {
	return Job{
		Id:         id,
		Status:     "Starting",
		Statuschan: c,
		Err:        nil,
		Created:    time.Now(),
		Modified:   time.Now(),
	}
}
