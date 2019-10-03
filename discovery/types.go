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
	// Instance is the unique identifier of a discoverer, e.g. the Hostname of the server where the discoverer is running.
	Instance string
	// Name is set by the Discoverer to a value that identifies the type of the Discoverer.
	Name string
	// Images is a list of discovered Docker images.
	Images []*Image
}

func (i *Input) String() string {
	return fmt.Sprintf("%s - %s", i.Instance, i.Name)
}

func ValidateImage(i *Image) []error {
	var errs []error
	if i.Digest == "" {
		errs = append(errs, fmt.Errorf("Field digest of Image '%s' is empty", i))
	}

	if i.Repository == "" {
		errs = append(errs, fmt.Errorf("Field repository of Image '%s' is empty", i))
	}

	if i.Source == "" {
		errs = append(errs, fmt.Errorf("Field source of Image '%s' is empty", i))
	}

	if i.Tag == "" {
		errs = append(errs, fmt.Errorf("Field tag of Image '%s' is empty", i))
	}

	return errs
}

func ValidateInput(i *Input) []error {
	var errs []error
	if i.Instance == "" {
		errs = append(errs, fmt.Errorf("Field instance of Input '%s' is empty", i))
	}

	if i.Name == "" {
		errs = append(errs, fmt.Errorf("Field name of Input '%s' is empty", i))
	}

	for _, img := range i.Images {
		errs = append(errs, ValidateImage(img)...)
	}

	return errs
}
