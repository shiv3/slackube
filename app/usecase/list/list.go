package list

import (
	"context"

	"github.com/shiv3/slackube/app/adapter/k8s/list"
)

type ListUseCaseImpl struct {
	ListAdapter list.ListAdapterInterface
}

func (u ListUseCaseImpl) ListNameSpace(ctx context.Context) ([]string, error) {
	nsList, err := u.ListAdapter.ListNs(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(nsList.Items))
	for _, item := range nsList.Items {
		ret = append(ret, item.Name)
	}
	return ret, err
}
