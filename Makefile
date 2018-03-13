PHONY: build build-docker

build: 
	dep ensure
	go build -o snapshot-app main.go

build-docker: build
	docker build -t snapshot .
