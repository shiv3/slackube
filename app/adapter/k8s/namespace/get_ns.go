package namespace

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type NameSpace struct {
	ClientSet *kubernetes.Clientset
}

func (ns *NameSpace) List() (*v1.NamespaceList, error) {
	return ns.ClientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
}
