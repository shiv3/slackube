package usecase

import (
	"github.com/shiv3/slackube/app/adapter/k8s"
	"github.com/shiv3/slackube/app/adapter/registory/gcr"
	gcr2 "github.com/shiv3/slackube/app/usecase/gcr"
	k8susecase "github.com/shiv3/slackube/app/usecase/k8s"
)

type UsecasesImpl struct {
	*k8susecase.K8sUseCaseImpl
	*gcr2.UseCaseImpl
}

func NewUsecasesImpl() (*UsecasesImpl, error) {
	client, err := k8s.NewK8SClientClient()
	if err != nil {
		return nil, err
	}
	gcrImpl := gcr.NewGcrAdapterImpl("asia.gcr.io")
	return &UsecasesImpl{
		K8sUseCaseImpl: &k8susecase.K8sUseCaseImpl{
			K8SAdapter: client,
		},
		UseCaseImpl: gcr2.NewUseCaseImpl(gcrImpl),
	}, nil
}
