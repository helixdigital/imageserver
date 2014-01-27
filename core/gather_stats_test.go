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
	"testing"

	"github.com/helixdigital/imageserver/entities"
	"github.com/helixdigital/imageserver/plugin/storage"
)

func TestTotalCount(t *testing.T) {
	jobstore := storage.NewJobStore()
	InjectStorageReporter(&jobstore)
	stored := TotalCount()
	if stored != 0 {
		t.Error("Newly-created JobStore should have total count of 0 but was", stored)
	}
	id := jobstore.AssignFreeId()
	c := make(<-chan entities.StatusMsg)
	job := entities.CreateJob(id, c)
	jobstore.AddJob(job)
	stored = TotalCount()
	if stored != 1 {
		t.Error("Newly-created JobStore should have total count of 0 but was", stored)
	}
	job.Status = "Done"
	jobstore.Replace(id, job)
	stored = jobstore.TotalCount()
	if stored != 1 {
		t.Error("Newly-created JobStore should still have total count of 1 but was", stored)
	}
}
