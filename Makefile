all:
	go build

run:
	./fetchstock ge

fmt:
	go fmt

packages:
	go get "github.com/cfdrake/go-ystocks"
