package discovery

import (
	"fmt"
	"time"
)

//Container is a discovered container.
type Container struct {
	// CreatedAt is the time at which the container was created.
	CreatedAt time.Time
	// Image is the image from which the container has been created.
	Image Image
	// Name is the name of the container.
	Name string
}

func (c *Container) String() string {
	return fmt.Sprintf("%s@%s", c.Name, c.CreatedAt)
}

// Image is the current image of a discovered container.
type Image struct {
	// Digest is the digest of the Docker image, e.g. "sha254:abc...".
	Digest string
	// Repository is the repository of the Docker image, e.g. "docker.io/imagespy/imagespy".
	Repository string
	// Tag is the tag of the Docker image, e.g. "latest".
	Tag string
}

func (i *Image) String() string {
	return fmt.Sprintf("%s:%s@%s", i.Repository, i.Tag, i.Digest)
}

// Result is the result returned by a Discoverer.
type Result struct {
	// Containers is a list of discovered Docker containers and their images.
	Containers []Container
}

// type ValidationError struct {
// 	errs []error
// }

// func (v *ValidationError) Error() string {
// 	b := bytes.NewBuffer([]byte{})
// 	for _, e := range v.errs {
// 		fmt.Fprintf(b, "%s\n", e)
// 	}

// 	return b.String()
// }

func ValidateContainer(c Container) (errs []error) {
	emptyT := time.Time{}
	if c.CreatedAt == emptyT {
		errs = append(errs, fmt.Errorf("Field createdAt of Container '%s' not set", c))
	}

	if c.Name == "" {
		errs = append(errs, fmt.Errorf("Field name of Container '%s' is empty", c))
	}

	errs = append(errs, ValidateImage(c.Image)...)
	return errs
}

func ValidateImage(i Image) []error {
	var errs []error
	if i.Digest == "" {
		errs = append(errs, fmt.Errorf("Field digest of Image '%s' is empty", i))
	}

	if i.Repository == "" {
		errs = append(errs, fmt.Errorf("Field repository of Image '%s' is empty", i))
	}

	if i.Tag == "" {
		errs = append(errs, fmt.Errorf("Field tag of Image '%s' is empty", i))
	}

	return errs
}

func ValidateResult(r *Result) []error {
	var errs []error
	for _, c := range r.Containers {
		errs = append(errs, ValidateContainer(c)...)
	}

	return errs
}
