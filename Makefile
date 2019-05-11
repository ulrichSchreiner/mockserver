.PHONY:
build:
	go build

.PHONY:
dockerbuild:
	docker build -t ulrichschreiner/mockserver .

.PHONY:
push:
	docker push ulrichschreiner/mockserver
