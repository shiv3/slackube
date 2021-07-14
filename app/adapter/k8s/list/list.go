package list

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type ListAdapterInterface interface {
	ListNs(ctx context.Context) (*corev1.NamespaceList, error)
	ListPod(ctx context.Context, nameSpace string) (*corev1.PodList, error)
	ListDeployment(ctx context.Context, nameSpace string) (*appsv1.DeploymentList, error)
	ListService(ctx context.Context, nameSpace string) (*corev1.ServiceList, error)
}

type ListAdapter struct {
	ClientSet *kubernetes.Clientset
}

func (a ListAdapter) ListNs(ctx context.Context) (*corev1.NamespaceList, error) {
	return a.ClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}

func (a ListAdapter) ListPod(ctx context.Context, nameSpace string) (*corev1.PodList, error) {
	return a.ClientSet.CoreV1().Pods(nameSpace).List(ctx, metav1.ListOptions{})
}

func (a ListAdapter) ListDeployment(ctx context.Context, nameSpace string) (*appsv1.DeploymentList, error) {
	return a.ClientSet.AppsV1().Deployments(nameSpace).List(ctx, metav1.ListOptions{})
}

func (a ListAdapter) ListService(ctx context.Context, nameSpace string) (*corev1.ServiceList, error) {
	return a.ClientSet.CoreV1().Services(nameSpace).List(ctx, metav1.ListOptions{})
}
