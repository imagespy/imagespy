VERSION ?= master
IMAGE = imagespy/imagespy

.PHONY: build
build:
	go build -o imagespy -mod=vendor core/cmd/main.go

.PHONY: package
package:
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: release
release: package
	docker push $(IMAGE):$(VERSION)
