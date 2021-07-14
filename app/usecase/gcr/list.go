package gcr

import (
	"context"
	"fmt"
	"regexp"

	"github.com/shiv3/slackube/app/adapter/registory/gcr"
)

type UseCaseImpl struct {
	GcrAdapter *gcr.GcrAdapterImpl
}

func NewUseCaseImpl(gcrAdapter *gcr.GcrAdapterImpl) *UseCaseImpl {
	return &UseCaseImpl{GcrAdapter: gcrAdapter}
}

func (u UseCaseImpl) ListImageTag(ctx context.Context, repoName string) ([]gcr.Image, error) {
	gcrRegex, _ := regexp.Compile("gcr.io")
	if gcrRegex.Match([]byte(repoName)) {
		return u.GcrAdapter.ListTags(repoName, "slackube")
	}

	return nil, fmt.Errorf("docker registory not supported")
}
