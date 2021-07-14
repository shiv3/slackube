package gcr

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/google/go-containerregistry/pkg/name"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

type GcrAdapterInterface interface {
	ListTags(ctx context.Context, fullname string, option *ListTagsOption) ([]Image, error)
}

type GcrAdapterImpl struct {
	host string
}

func NewGcrAdapterImpl(host string) *GcrAdapterImpl {
	return &GcrAdapterImpl{host: host}
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

func (s GcrAdapterImpl) ListTags(imageName, userAgent string) ([]Image, error) {
	repo, err := name.NewRepository(imageName, name.WeakValidation)
	if err != nil {
		return nil, err
	}
	tags, err := google.List(repo,
		google.WithAuthFromKeychain(authn.DefaultKeychain),
		google.WithTransport(http.DefaultTransport),
		google.WithUserAgent(userAgent),
		google.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}
	res := make([]Image, 0, len(tags.Manifests))
	for k, v := range tags.Manifests {
		res = append(res, Image{
			Digest:   k,
			Tags:     v.Tags,
			Uploaded: v.Uploaded,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[j].Uploaded.Sub(res[i].Uploaded) < 0
	})
	return res, err
}
