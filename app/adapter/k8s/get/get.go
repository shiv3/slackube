package get

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type GetAdapterInterface interface {
	GetNs(ctx context.Context, name string) (*corev1.Namespace, error)
	GetPod(ctx context.Context, name string, nameSpace string) (*corev1.Pod, error)
	GetDeployment(ctx context.Context, name string, nameSpace string) (*appsv1.Deployment, error)
	GetService(ctx context.Context, name string, nameSpace string) (*corev1.Service, error)
}

type GetAdapter struct {
	ClientSet *kubernetes.Clientset
}

func (a GetAdapter) GetNs(ctx context.Context, name string) (*corev1.Namespace, error) {
	return a.ClientSet.CoreV1().Namespaces().Get(ctx, name, v1.GetOptions{})
}

func (a GetAdapter) GetPod(ctx context.Context, nameSpace string, name string) (*corev1.Pod, error) {
	return a.ClientSet.CoreV1().Pods(nameSpace).Get(ctx, name, metav1.GetOptions{})
}

func (a GetAdapter) GetDeployment(ctx context.Context, nameSpace string, name string) (*appsv1.Deployment, error) {
	return a.ClientSet.AppsV1().Deployments(nameSpace).Get(ctx, name, metav1.GetOptions{})
}

func (a GetAdapter) GetService(ctx context.Context, nameSpace string, name string) (*corev1.Service, error) {
	return a.ClientSet.CoreV1().Services(nameSpace).Get(ctx, name, metav1.GetOptions{})
}
