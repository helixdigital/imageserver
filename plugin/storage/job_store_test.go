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

package storage

import (
	"testing"

	"github.com/helixdigital/imageserver/entities"
)

func TestAssignFreeId(t *testing.T) {
	jobstore := NewJobStore()
	result := jobstore.AssignFreeId()
	if result != 0 {
		t.Error("AssignFreeId should have returned 0 but returned", result)
	}
	result = jobstore.AssignFreeId()
	if result != 1 {
		t.Error("AssignFreeId should have returned 1 but returned", result)
	}
	result = jobstore.AssignFreeId()
	if result != 2 {
		t.Error("AssignFreeId should have returned 2 but returned", result)
	}
}

func TestAdd(t *testing.T) {
	jobstore := NewJobStore()
	id := jobstore.AssignFreeId()

	_, ok := jobstore.GetJob(id)
	if ok {
		t.Error("Getting an ID before adding something to it should not have been OK")
	}

	c := make(<-chan entities.StatusMsg)
	job := entities.CreateJob(id, c)
	jobstore.AddJob(job)

	outjob, ok := jobstore.GetJob(id)
	if !ok {
		t.Error("Getting an ID after adding something to it should have been OK")
	}

	if outjob != job {
		t.Error("Was supposed to get the same job from jobstore but didn't")
	}
}

func TestEdit(t *testing.T) {
	jobstore := NewJobStore()
	id := jobstore.AssignFreeId()
	c := make(<-chan entities.StatusMsg)
	job := entities.CreateJob(id, c)
	jobstore.AddJob(job)
	outjob, _ := jobstore.GetJob(id)
	if outjob.Status != "Starting" {
		t.Error("Status before editing should have been 'Starting' but was", outjob.Status)
	}
	outjob.Status = "Done"
	jobstore.Replace(id, outjob)
	editedjob, _ := jobstore.GetJob(id)
	if editedjob.Status != "Done" {
		t.Error("Status after editing should have been 'Done' but was", outjob.Status)
	}
}

func TestTotalCount(t *testing.T) {
	jobstore := NewJobStore()
	stored := jobstore.TotalCount()
	if stored != 0 {
		t.Error("Newly-created JobStore should have total count of 0 but was", stored)
	}
	id := addOneJob(&jobstore)
	stored = jobstore.TotalCount()
	if stored != 1 {
		t.Error("Newly-created JobStore should have total count of 1 but was", stored)
	}
	editJob(&jobstore, id)
	stored = jobstore.TotalCount()
	if stored != 1 {
		t.Error("Newly-created JobStore should still have total count of 1 but was", stored)
	}
}

func TestCountByStatus(t *testing.T) {
	jobstore := NewJobStore()
	by_status := jobstore.CountByStatus()
	if len(by_status) != 0 {
		t.Error("Newly-created JobStore should have total size of 0 but was", by_status)
	}
	id := addOneJob(&jobstore)
	by_status = jobstore.CountByStatus()
	if len(by_status) != 1 {
		t.Error("Newly-created JobStore should have total size of 1 but was", by_status)
	}
	if by_status["Starting"] != 1 {
		t.Error("Should have 1 starting status but was", by_status["Starting"])
	}
	addOneJob(&jobstore)
	by_status = jobstore.CountByStatus()
	if len(by_status) != 1 {
		t.Error("Newly-created JobStore should have total size of 1 but was", by_status)
	}
	if by_status["Starting"] != 2 {
		t.Error("Should have 2 starting status but was", by_status["Starting"])
	}
	if by_status["Done"] != 0 {
		t.Error("Should have 0 done status but was", by_status["Done"])
	}
	editJob(&jobstore, id)
	by_status = jobstore.CountByStatus()
	if len(by_status) != 2 {
		t.Error("Newly-created JobStore should still have total size of 2 but was", by_status)
	}
	if by_status["Starting"] != 1 {
		t.Error("Should have 1 starting status but was", by_status["Starting"])
	}
	if by_status["Done"] != 1 {
		t.Error("Should have 1 done status but was", by_status["Done"])
	}
}

func addOneJob(jobstore entities.JobStore) int {
	id := jobstore.AssignFreeId()
	c := make(<-chan entities.StatusMsg)
	job := entities.CreateJob(id, c)
	jobstore.AddJob(job)
	return id
}

func editJob(jobstore entities.JobStore, id int) {
	job, _ := jobstore.GetJob(id)
	job.Status = "Done"
	jobstore.Replace(id, job)
}
