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

// Getting data from the user into the use cases and back out again
// is not a core functionality of an application. It is a detail, a
// plugin. This package displays pages to the user any way at all.
// Currently it is with basic POSTs and GETs with almost no structure
// or protocol. This will hopefully change but the interface between
// the user and the use-case code in any application is volatile -
// likely to change rapidly - and should be well-separated from other
// code. This package is for that purpose.
//
// Most handlers in this package will be taking data from the http
// call and massaging it into simple, basic datastructures that are
// passed into use-cases found in the core package and then getting
// simple data structures back from the use-case and displaying them
// in an appropriate way to the user.
package presentation

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/helixdigital/imageserver/core"
)

// Start the go standard library web server on the given port.
// It sets up handlers for various URLs that are private to this package.
func StartWebServer(port int) {
	setuphandlers()
	portstring := fmt.Sprintf(":%d", port)
	fmt.Println("Running on ", portstring)
	log.Fatal(http.ListenAndServe(portstring, nil))
}

// There are four endpoints:
// - `/` Does nothing at the moment: merely displays a hello world
// - `/status` returns current status of the given job
// - `/stats` returns the current status of the running server
// - `/request` starts a new job
func setuphandlers() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/request", requestHandler)
}

// Simply returns a greeting string when `/` is called
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Image Server", http.StatusNotImplemented)
}

// Calls the core.JobStatus use-case with the jobid found in the GET
// query parmeter.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	jobid := toInt(r.FormValue("jobid"))
	status, err := core.JobStatus(jobid)
	if strings.Contains(status, "No job found with id") {
		fmt.Printf("No job found with id %d\n", jobid)
		http.Error(w, status, http.StatusGone)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	}
	fmt.Printf("Status %d is %s\n", jobid, status)
	fmt.Fprintf(w, "%s", status)
}

// Calls the core.NewJob use-case with the data send in the http POST form
func requestHandler(w http.ResponseWriter, r *http.Request) {
	jobreq := getJobRequestFrom(r)
	if r.FormValue("debug") == "1" {
		fmt.Fprintf(w, "%#v", jobreq)
		return
	}

	newid := core.NewJob(jobreq)
	fmt.Printf("New job requested %#v -> id:%d\n", newid)
	fmt.Fprintf(w, "%d", newid)

}

// Displays as JSON the structure returned by the call to core.GetStats
func statsHandler(w http.ResponseWriter, r *http.Request) {
	data := core.GetStats()
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "%s", b)
}

// Converts the data sent in the http POST to requestHandler into a
// core.JobRequest data structure.
func getJobRequestFrom(r *http.Request) core.JobRequest {
	return core.JobRequest{
		Local_filename: r.FormValue("local_filename"),
		Crop_to: image.Rect(
			toInt(r.FormValue("crop_to_x")),
			toInt(r.FormValue("crop_to_y")),
			intAndAdd(r, "crop_to_x", "crop_to_w"),
			intAndAdd(r, "crop_to_y", "crop_to_h"),
		),
		Resize_width:      toUint(r.FormValue("resize_width")),
		Resize_height:     toUint(r.FormValue("resize_height")),
		Uploaded_filename: r.FormValue("uploaded_filename"),
	}
}

// When passed two keys to the form values in r, this gets the
// corresponding values - which will be strings - converts these
// to ints and then returns the sum of these two ints.
//
// TODO: Methods should do one thing only. Refactor this.
func intAndAdd(r *http.Request, alice string, bob string) int {
	return toInt(r.FormValue(alice)) + toInt(r.FormValue(bob))
}

// TODO: Refactor this duplication
func toUint(input string) uint {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}
func toInt(input string) int {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0
	}
	return int(i)
}
