package discovery

import (
	"fmt"
)

type Image struct {
	Digest     string
	Repository string
	Source     string
	Tag        string
}

func (i *Image) String() string {
	return fmt.Sprintf("%s - %s:%s@%s", i.Source, i.Repository, i.Tag, i.Digest)
}

type Input struct {
	Name   string
	Images []*Image
}
