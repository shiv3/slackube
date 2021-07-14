package gcr

import (
	"context"

	"github.com/shiv3/slackube/app/adapter/registory/gcr"
)

type UseCaseImpl struct {
	GcrAdapter *gcr.GcrAdapterImpl
}

func NewUseCaseImpl(gcrAdapter *gcr.GcrAdapterImpl) *UseCaseImpl {
	return &UseCaseImpl{GcrAdapter: gcrAdapter}
}

func (u UseCaseImpl) ListImageTag(ctx context.Context, repoName string) ([]gcr.Image, error) {
	return u.GcrAdapter.ListTags(repoName, "slackube")
}
