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

.PHONY: test
test:
	go test -mod=vendor ./...

.PHONY: bootstrap_test
bootstrap_test:
	docker-compose -f core/fixtures/TestRunner_Run/registry/docker-compose.yml -p testrunner up -d
	docker pull redis@sha256:e59e6cab7ada6fa8cfeb9ad4c5b82bf7947bf3620cafc17e19fdbff01b239981
	docker tag redis@sha256:e59e6cab7ada6fa8cfeb9ad4c5b82bf7947bf3620cafc17e19fdbff01b239981 127.0.0.1:52854/redis:4.0.14-alpine
	docker push 127.0.0.1:52854/redis:4.0.14-alpine
	docker pull redis@sha256:d9ea76b14d4771c7cd0c199de603f3d9b1ea246c0cbaae02b86783e1c1dcc3d1
	docker tag redis@sha256:d9ea76b14d4771c7cd0c199de603f3d9b1ea246c0cbaae02b86783e1c1dcc3d1 127.0.0.1:52854/redis:5.0.6-alpine
	docker push 127.0.0.1:52854/redis:5.0.6-alpine
