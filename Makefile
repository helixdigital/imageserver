GOPATH = $(shell cd $(CURDIR)/../../../..; pwd)
current: tag test

build:
	go build .

tag:
	ctags --recurse=yes .

fullbuild:
	go get github.com/nfnt/resize
	go get launchpad.net/goamz/aws
	go get launchpad.net/goamz/s3
	go build .

test:
	go test ./...

cover:
	go test -cover ./...

coverage:
	go test -coverprofile=c.out github.com/helixdigital/imageserver/core && go tool cover -func=c.out 
	echo ""
	go test -coverprofile=c.out github.com/helixdigital/imageserver/plugin/presentation && go tool cover -func=c.out 
	echo ""
	go test -coverprofile=c.out github.com/helixdigital/imageserver/plugin/storage && go tool cover -func=c.out 
	echo ""
	go test -coverprofile=c.out github.com/helixdigital/imageserver/plugin/upload && go tool cover -func=c.out 
	rm c.out
