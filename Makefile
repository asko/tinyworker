GOSRC := $(shell find . -name '*.go' ! -name '*.pb.go' ! -name '*test*' ! -path './vendor/*' ! -path './cmd/*')
build=@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s -w"

clean:
	rm -rf tinyworker

build: main.go $(GOSRC)
	go build -o tinyworker $<

build-static: main.go $(GOSRC)
	${build} -o tinyworker $<

.PHONY: test build build-static .force
