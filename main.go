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

package main

import (
	"flag"
	"os"

	"github.com/helixdigital/imageserver/core"
	"github.com/helixdigital/imageserver/plugin/presentation"
	"github.com/helixdigital/imageserver/plugin/storage"
	"github.com/helixdigital/imageserver/plugin/upload"
)

func injectDependencies() {
	uploader := upload.NewAmazonS3Upload(s3accesskey, s3secretkey, s3bucketname)
	core.InjectUploader(uploader)

	store := storage.NewJobStore()
	core.InjectJobstore(&store)
	core.InjectStorageReporter(&store)
}

var portflag int
var s3accesskey string
var s3secretkey string
var s3bucketname string

func handleFlags() {
	flag.IntVar(&portflag, "port", 9877, "The port the app will run on")
	flag.StringVar(
		&s3accesskey,
		"accesskey",
		os.Getenv("IMAGESERVER_S3_ACCESS_KEY"),
		"Amazon S3 access key",
	)
	flag.StringVar(
		&s3secretkey,
		"secretkey",
		os.Getenv("IMAGESERVER_S3_SECRET_KEY"),
		"Amazon S3 secret key",
	)
	flag.StringVar(
		&s3bucketname,
		"bucketname",
		os.Getenv("IMAGESERVER_S3_BUCKET_NAME"),
		"Amazon S3 bucket name",
	)
	flag.Parse()
}

func main() {
	handleFlags()
	injectDependencies()
	presentation.StartWebServer(portflag)
}
