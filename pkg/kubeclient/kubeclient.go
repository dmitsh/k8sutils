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
}

func New() (*Client, error) {
	client, config, err := createKubeClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
		config: config,
	}, nil
}

func createKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	var kubeconfig string

	config, err := rest.InClusterConfig()
	if err != nil {
		// try outcluster config
		kubeconfig = os.Getenv("KUBECONFIG")
		if len(kubeconfig) == 0 {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			log.Warnf("$KUBECONFIG is not set. Trying %s", kubeconfig)
			if _, err := os.Stat(kubeconfig); err != nil {
				return nil, nil, fmt.Errorf("failed to find kubeconfig")
			}
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, err
		}
	}
	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return client, config, nil
}
