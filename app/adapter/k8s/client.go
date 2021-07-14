package k8s

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/rest"

	"github.com/shiv3/slackube/app/adapter/k8s/img"

	"github.com/shiv3/slackube/app/adapter/k8s/get"

	"github.com/shiv3/slackube/app/adapter/k8s/list"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SAdapter struct {
	List  list.ListAdapterInterface
	Get   get.GetAdapterInterface
	Image img.ImageAdapterInterface
}

func NewK8SClientClient() (K8SAdapter, error) {
	// first use InClusterConfig
	config, err := rest.InClusterConfig()
	if err != nil {
		// if couldn't get InClusterConfig, try to use local kubeconfig file
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return K8SAdapter{}, err
		}
	}
	ks, err := kubernetes.NewForConfig(config)
	if err != nil {
		return K8SAdapter{}, err
	}
	// create the clientset
	return K8SAdapter{
		List:  list.ListAdapter{ClientSet: ks},
		Get:   get.GetAdapter{ClientSet: ks},
		Image: img.ImageAdapter{ClientSet: ks},
	}, err

}
