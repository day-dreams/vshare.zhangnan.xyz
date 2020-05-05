.PHONY: build run


WORKDIR=$(CURDIR)

build:
	go build -o build/VShareServer main.go
run: build
	build/VShareServer \
	--vfile=$(WORKDIR)/runtime/config.json \
	--vindex=$(WORKDIR)/runtime/index.html
daemon: build
	(pkill VShareServer || echo "skip kill" )&& nohup build/VShareServer \
	--vfile=$(WORKDIR)/runtime/config.json \
	--vindex=$(WORKDIR)/runtime/index.html \
	&> ./build/vshare.log &
buildLinux:
	GOOS=linux GOARCH=amd64 go build -o build/VShareServer main.go
docker: buildLinux
	docker run --rm -it `docker build -q .`
