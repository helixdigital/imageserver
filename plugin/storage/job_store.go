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

// Persistent storage, like presentation, is not a core functionality. It is
// fragile to change and should have as few things depend on it as possible.
package storage

import (
	"sync"
	"time"

	"github.com/helixdigital/imageserver/entities"
)

// Implements both entities.JobStore and core.StorageReporter
type jobs struct {
	store   map[int]entities.Job
	id_lock sync.Mutex
	next_id int
}

func (self *jobs) AddJob(newjob entities.Job) {
	(*self).store[newjob.Id] = newjob
}

// Update existing job
func (self *jobs) Replace(id int, newversion entities.Job) {
	newversion.Modified = time.Now()
	(*self).store[id] = newversion
}

func (self *jobs) AssignFreeId() (next int) {
	(*self).id_lock.Lock()
	defer (*self).id_lock.Unlock()
	next = (*self).next_id
	(*self).next_id = (*self).next_id + 1
	return
}

func (self *jobs) GetJob(id int) (entities.Job, bool) {
	job, ok := (*self).store[id]
	return job, ok
}

// NewJobStore is a factory for an empty collection of jobs
func NewJobStore() jobs {
	return jobs{store: make(map[int]entities.Job, 0)}
}

func (self *jobs) TotalCount() int {
	return len((*self).store)
}

func (self *jobs) CountByStatus() map[string]int {
	output := make(map[string]int)
	for _, job := range self.store {
		key := job.Status
		count, ok := output[key]
		if !ok {
			count = 0
		}
		output[key] = count + 1
	}
	return output
}
