package kubecli

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeCli wraps k8s.io clientSet.
type KubeCli struct {
	Config        *rest.Config
	ClientSet     *kubernetes.Clientset
	DynamicClient dynamic.Interface
	Namespace     string
}

// KubeCliOpts declare options to customize the generated client.
type KubeCliOpts struct {
	// QPS indicates the maximum QPS to the master from this client.
	// If it's zero, the created RESTClient will use DefaultQPS: 5
	QPS float32

	// Maximum burst for throttle.
	// If it's zero, the created RESTClient will use DefaultBurst: 10.
	Burst int
}

// New creates a new KubeCli client from Kubeconfig.
func New(context string) (*KubeCli, error) {
	return NewWithOptions(context, nil)
}

// NewWithOptions creates a new KubeCli client from Kubeconfig
// with come configuration options.
func NewWithOptions(context string, opts *KubeCliOpts) (*KubeCli, error) {
	var config *rest.Config
	var namespace string

	config, namespace, err := outOfClusterConfig(context)
	if err != nil {
		// Maybe we are "in cluster", let's try this
		config, err = inClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	if opts != nil {
		config.QPS = opts.QPS
		config.Burst = opts.Burst
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create ClientSet: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create dynamicClient: %w", err)
	}

	return &KubeCli{
		Config:        config,
		ClientSet:     clientSet,
		Namespace:     namespace,
		DynamicClient: dynamicClient,
	}, nil
}

// NewOutOfCluster creates a KubeCli instance using only "out-of-cluster" configuration.
func NewOutOfCluster(context string) (*KubeCli, error) {
	return NewOutOfClusterWithOptions(context, nil)
}

// NewOutOfClusterWithOptions creates a KubeCli instance using only
// "out-of-cluster" configuration with come configuration options.
func NewOutOfClusterWithOptions(context string, opts *KubeCliOpts) (*KubeCli, error) {
	config, namespace, err := outOfClusterConfig(context)
	if err != nil {
		return nil, fmt.Errorf("could not create Kube Configuration, reason: %w", err)
	}

	if opts != nil {
		config.QPS = opts.QPS
		config.Burst = opts.Burst
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create ClientSet: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create dynamicClient: %w", err)
	}

	return &KubeCli{
		Config:        config,
		ClientSet:     clientSet,
		Namespace:     namespace,
		DynamicClient: dynamicClient,
	}, nil
}

func outOfClusterConfig(context string) (*rest.Config, string, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	override := &clientcmd.ConfigOverrides{}
	if context != "" {
		override.CurrentContext = context
	}

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}, override)
	clientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("could not create Kube Configuration outOfCluster, possible reasons: %w", err)
	}

	namespace, _, err := config.Namespace()
	if err != nil {
		return nil, "", fmt.Errorf("could not create Kube Configuration outOfCluster, possible reasons: %w", err)
	}

	return clientConfig, namespace, nil
}

func inClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("could not create Kube Configuration inCluster, possible reasons: %w", err)
	}
	return config, nil
}
