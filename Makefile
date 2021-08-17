SHELL = bash
PROJECT := terraform-provider-krok
VERSION ?= v0.0.1

all: binaries

clean:
	rm -Rf bin

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64 darwin/arm64" \
		-ldflags="-X main.version=${VERSION}" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)" \
		-tags="netgo" \
		./...

bootstrap:
	go get github.com/hashicorp/terraform-plugin-sdk/plugin
	go get github.com/hashicorp/terraform-plugin-sdk/terraform

docker:
	docker build \
		--build-arg=GOARCH=amd64 \
		-t $(DOCKERIMAGE) .

docker-push:
	docker push $(DOCKERIMAGE)

