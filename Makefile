.PHONY: build
build:
	go build -o imagespy -mod=vendor core/cmd/main.go

package:
	docker build -t imagespy/imagespy .
