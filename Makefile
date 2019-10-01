VERSION ?= master
IMAGE = quay.io/imagespy/imagespy

.PHONY: build
build:
	go build -o imagespy -mod=vendor core/cmd/main.go

package:
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: release
release: package
	docker push $(IMAGE):$(VERSION)
