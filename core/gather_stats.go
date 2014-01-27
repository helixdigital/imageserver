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

import "time"

// StorageReporter is the plugin that provides statistics
// of the database
type StorageReporter interface {
	TotalCount() int
	CountByStatus() map[string]int
}

var reporter StorageReporter

// InjectStorageReporter is the setter for the current StorageReporter
func InjectStorageReporter(storagereporter StorageReporter) {
	reporter = storagereporter
}

// TotalCount gives the number of jobs run on the server since
// it started.
func TotalCount() int {
	return reporter.TotalCount()
}

// CountByStatus returns the number of jobs partitioned by
// their current status
func CountByStatus() map[string]int {
	return reporter.CountByStatus()
}

var starttime = time.Now()

// SecondsUp is the difference between the current time and
// the time this server was started
func SecondsUp() int {
	return int(time.Now().Sub(starttime) / time.Second)
}

// Stats is the output data structure of the GetStats() function
type Stats struct {
	SecondsUp     int
	TotalCount    int
	CountByStatus map[string]int
}

// GetStats returns some basic statistics
// of the current state of this server
func GetStats() Stats {
	return Stats{SecondsUp(), TotalCount(), CountByStatus()}
}
