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

package upload

import (
	"log"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

// AmazonS3Upload implements github.com/helixdigital/imageserver/core/Uploader
//
// It
type AmazonS3Upload struct {
	accesskey  string
	secretkey  string
	bucketname string
}

// NewAmazonS3Upload is a factory that creates  an uploader.
// The parameters describe the account and bucket to upload files to.
// Note that the bucketname is set. To upload to different buckets, create
// new AmazonS3Uploads for each one.
func NewAmazonS3Upload(accesskey string, secretkey string, bucketname string) AmazonS3Upload {
	if accesskey == "" || secretkey == "" || bucketname == "" {
		log.Fatal(`
        Amazon credentials must be declared.
        Use environment variables:
        IMAGESERVER_S3_ACCESS_KEY
        IMAGESERVER_S3_SECRET_KEY
        IMAGESERVER_S3_BUCKET_NAME
        `)
	}
	return AmazonS3Upload{accesskey, secretkey, bucketname}
}

// Upload implements github.com/helixdigital/imageserver/core/Uploader interface.
// Saves the data in the io.Reader parameter on Amazon S3
func (self AmazonS3Upload) Upload(data []byte, mime string, uplname string) error {
	bucket := self.getbucket()
	return bucket.Put(uplname, data, mime, s3.PublicRead)
}

// Delete removes the path from Amazon S3.
func (self AmazonS3Upload) Delete(uplname string) error {
	bucket := self.getbucket()
	return bucket.Del(uplname)
}

func (self AmazonS3Upload) getbucket() *s3.Bucket {
	auth := aws.Auth{self.accesskey, self.secretkey}
	s := s3.New(auth, aws.APSoutheast2)
	return s.Bucket(self.bucketname)
}
