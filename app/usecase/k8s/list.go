package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/shiv3/slackube/app/adapter/k8s"
)

type K8sUseCaseImpl struct {
	K8SAdapter k8s.K8SAdapter
}

type ListNameSpaceOut struct {
	Name   string
	Status string
	Age    string
}

func (u K8sUseCaseImpl) GetNameSpace(ctx context.Context) ([]ListNameSpaceOut, error) {
	nsList, err := u.K8SAdapter.List.ListNs(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]ListNameSpaceOut, 0, len(nsList.Items))
	for _, item := range nsList.Items {
		res = append(res, ListNameSpaceOut{
			Name:   item.Name,
			Status: string(item.Status.Phase),
			Age:    getAge(item.CreationTimestamp.Time),
		})
	}
	return res, nil
}

func getAge(t time.Time) string {
	sub := time.Now().Sub(t).Hours()
	if sub > 24 {
		return fmt.Sprintf("%.0fd", sub/24)
	}
	return fmt.Sprintf("%.0fh", sub)
}

func (u K8sUseCaseImpl) ListNameSpace(ctx context.Context) ([]string, error) {
	nsList, err := u.K8SAdapter.List.ListNs(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(nsList.Items))
	for _, item := range nsList.Items {
		res = append(res, item.Name)
	}

	return res, nil
}

type DeploymentContainers struct {
	NameSpace       string
	Deployment      string
	ContainerImages map[string]string
}

func (u K8sUseCaseImpl) ListDeployment(ctx context.Context, ns string) ([]DeploymentContainers, error) {
	list, err := u.K8SAdapter.List.ListDeployment(ctx, ns)
	if err != nil {
		return nil, err
	}

	res := make([]DeploymentContainers, 0, len(list.Items))
	for _, item := range list.Items {
		c := make(map[string]string)
		for _, container := range item.Spec.Template.Spec.Containers {
			c[container.Name] = container.Image
		}
		res = append(res, DeploymentContainers{
			NameSpace:       ns,
			Deployment:      item.Name,
			ContainerImages: c,
		})
	}

	return res, nil
}

func (u K8sUseCaseImpl) ListContainer(ctx context.Context, name, ns string) ([]string, error) {
	dp, err := u.K8SAdapter.Get.GetDeployment(ctx, name, ns)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(dp.Spec.Template.Spec.Containers))
	for _, container := range dp.Spec.Template.Spec.Containers {
		res = append(res, container.Image)
	}
	return res, nil
}

func (u K8sUseCaseImpl) UpdateImage(ctx context.Context, ns, deployment, container, digest string) ([]string, error) {
	dp, err := u.K8SAdapter.Image.UpdateImage(ctx, ns, deployment, container, digest)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(dp.Spec.Template.Spec.Containers))
	for _, container := range dp.Spec.Template.Spec.Containers {
		res = append(res, container.Image)
	}
	return res, nil
}
