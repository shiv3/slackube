package dockerhub

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/heroku/docker-registry-client/registry"
)

type ListImageAdapterInterface interface {
	ListTags(ctx context.Context, fullname string, option *ListTagsOption) ([]Image, error)
}

type ListImageAdapterImpl struct {
	host string
}

func NewListImageAdapterImpl(host string) *ListImageAdapterImpl {
	return &ListImageAdapterImpl{host: host}
}

type ListTagsOption struct {
	getImageLimit int
}

// Image represents a image object of gks
type Image struct {
	Digest   string   `json:"digest",yaml:"metadata,flow"`
	Tags     []string `json:"tags,omitempty",yaml:"spec,flow"`
	Uploaded time.Time
}

func (s ListImageAdapterImpl) ListTags(imageName, userAgent string) ([]Image, error) {
	host := strings.Split(imageName, "/")[0]
	hub, err := registry.New(host, "", "")
	tags, err := hub.Tags(imageName)
	if err != nil {
		return nil, err
	}
	res := make([]Image, 0, len(tags))
	sort.Slice(res, func(i, j int) bool {
		return res[j].Uploaded.Sub(res[i].Uploaded) < 0
	})
	return res, err
}
