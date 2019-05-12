export GO111MODULE := on

.PHONY:
build:
	go build

.PHONY:
dockerbuild:
	docker build -t ulrichschreiner/mockserver:latest .

.PHONY:
push:
	docker push ulrichschreiner/mockserver:latest

.PHONY:
localtestservice:
	go run *.go -debug=true testfiles/usermocks.yaml
