package k8s

import (
	"flag"
	"path/filepath"

	"github.com/shiv3/slackube/app/adapter/k8s/list"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SClient struct {
	list.ListAdapterInterface
}

func NewK8SClientClient() (K8SClient, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	ks, err := kubernetes.NewForConfig(config)

	// create the clientset
	return K8SClient{
		ListAdapterInterface: list.ListAdapter{ClientSet: ks},
	}, err

}
