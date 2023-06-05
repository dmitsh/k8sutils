package kubeclient

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	client *kubernetes.Clientset
	config *rest.Config
	//mux    sync.Mutex
}

func New() (*Client, error) {
	client, config, err := createKubeClient()
	if err != nil {
		log.Errorf("could not create kubernetes client: %v", err)
		return nil, err
	}
	return &Client{
		client: client,
		config: config,
	}, nil
}

func createKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	var kubeconfig string
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			log.Warnf("$KUBECONFIG is not set. Using $HOME/.kube/config (if present)")
			kubeconfig = filepath.Join(
				os.Getenv("HOME"), ".kube", "config",
			)
		}
		if kubeconfig == "" {
			msg := "could not find kubeconfig"
			return nil, nil, fmt.Errorf(msg)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, err
		}
	}
	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	return client, config, nil
}
