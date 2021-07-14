package img

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type ImageAdapterInterface interface {
	UpdateImage(ctx context.Context, nameSpace, deployment, container, digest string) (*v1.Deployment, error)
}

type ImageAdapter struct {
	ClientSet *kubernetes.Clientset
}

func (a ImageAdapter) UpdateImage(ctx context.Context, nameSpace, deployment, container, digest string) (*v1.Deployment, error) {
	dp, err := a.ClientSet.AppsV1().Deployments(nameSpace).Get(ctx, deployment, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	for i, c := range dp.Spec.Template.Spec.Containers {
		if c.Name == container {
			dp.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s@%s", getImageName(c.Image), digest)
		}
	}
	return a.ClientSet.AppsV1().Deployments(nameSpace).Update(ctx, dp, metav1.UpdateOptions{})
}

func getImageName(imageName string) string {
	image := strings.Split(strings.Split(imageName, ":")[0], "@")[0]
	return image
}
