package discovery

import (
	"fmt"
)

// Image is a discovered Docker image.
type Image struct {
	// Digest is the digest of the Docker image, e.g. "sha254:abc...".
	Digest string
	// Repository is the repository of the Docker image, e.g. "docker.io/imagespy/imagespy".
	Repository string
	// Source is set by the Discoverer to a value that uniquely identifies the source of the image, e.g. the name of the Docker container.
	Source string
	// Tag is the tag of the Docker image, e.g. "latest".
	Tag string
}

func (i *Image) String() string {
	return fmt.Sprintf("%s - %s:%s@%s", i.Source, i.Repository, i.Tag, i.Digest)
}

// Input is the result returned by a Discoverer.
type Input struct {
	// Name is set by the Discoverer to a value that uniquely identifies the Discoverer.
	Name string
	// Images is a list of discovered Docker images.
	Images []*Image
}
