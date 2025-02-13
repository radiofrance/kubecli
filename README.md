# KubeCli

_KubeCli_ is a simplified golang client for the Kubernetes API. [It is built around the official golang Kubernetes client](https://godoc.org/k8s.io/client-go/kubernetes).

# Quick Start

```go
import "github.com/radiofrance/kubecli"

// Create a new instance of KubeCli for a specified context (leave empty for default context)
// KubeCli will automatically select an "In-Cluster" or "Out-of-Cluster" configuration method
kube, err := kubecli.New("myContext")
if err != nil {
    return fmt.Errorf("could not create kubecli: %v", err)
}

// You can use built-in KubeCli methods
replicaSet, err := kube.GetReplicaSetByName("myDeploymentName")
if err != nil {
    return fmt.Errorf("could not get active ReplicaSet for deployment, reason: %v", err)
}

pods, err := kube.FindPods(replicaSet)
if err != nil {
    return fmt.Errorf("could not find Pods in replicaset, reason: %v", err)
}

// Or directly access k8s.io/client-go/kubernetes.ClientSet
deployment, err := kube.ClientSet.AppsV1().Deployments(metav1.NamespaceDefault).Get("myDeploymentName", metav1.GetOptions{})
if err != nil {
    return fmt.Errorf("could not get Deployment: %v", err)
}

``
