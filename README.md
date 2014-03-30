Imageserver
===========

[v.1.0.1](http://semver.org/)



What is it?
-----------

Crop, resize and push to S3 your user's images.
Originally written as a reusable webapp component.
Install it on the same server as your webapp.

How to use it?
--------------

`imageserver --port=9877 [--s3accesskey="c0ffee"] [--s3secretkey="cafe"] [--s3bucketname="mybucket"]`

`--port` The port the imageserver will serve

if any of the --s3... parameters are missing, they must be specified in the environment variables:
`IMAGESERVER_S3_ACCESS_KEY`
`IMAGESERVER_S3_SECRET_KEY`
`IMAGESERVER_S3_BUCKET_NAME`


A new crop, resize and upload job is created by POSTing to `/request` with the following eight form elements:
`local_filename` (the name of the file on the local filesystem to use as input)
`crop_to_x, crop_to_y, crop_to_w, crop_to_h` (the rectangle of the input image that will be visible in the end)
`resize_width, resize_height` (the dimensions of the final image after the cropped image is resized - if one of these is "0" then the other resize parameter is used to size the image with aspect preserved. Both can be "0" in which case the image will not be resized)
`uploaded_filename` (the name that the resized image will be stored on S3 as)

The POST to `/request` will return a body with a single string as response. This string is the `jobid`.

Subsequently GETting from `/status?jobid=[jobid]` (that is, with a GET query that has a key of `jobid` and a value being the string returned from the original POST to `/request`) will return in the body of the response only a single string that will be one of:
* "Reading the file"
* "Cropping"
* "Resizing"
* "Uploading"
* "Done"
* "Error reading the file"
* "Error in cropping"
* "Error in resizing"
* "Error in uploading"
* "Timed out"

The wording may change in the future. More may be added, Some of these may be removed.

Calls to `/stats` returns a JSON data structure showing a couple of rudimentary statistics describing the state of the server. The content of the response may change in the future. 


What are its limitations?
-------------------------

You do not pass in an image file to this service, only the name of the imagefile so you must run it on the same server as the webapp that accepts the file upload POSTs from the user.

It does no security, authentication, or authorisation. You are to protect this with your firewall.

The current implementation will collect jobs in a data structure in memory and they never get purged so it will continue to eat memory. It may need to be restarted every couple of months or so.

How do I compile it?
--------------------

* `go get github.com/helixdigital.com/imageserver`
* `make fullcompile`

How is it licensed?
-------------------

[Apache Software License 2](http://www.apache.org/licenses/LICENSE-2.0.html)
