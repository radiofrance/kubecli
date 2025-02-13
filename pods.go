package kubecli

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindPods returns the pod associated with a replicaset.
func (k *KubeCli) FindPods(replicaset *appsv1.ReplicaSet) ([]corev1.Pod, error) {
	pods, err := k.ClientSet.CoreV1().Pods(replicaset.GetNamespace()).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list pods: %w", err)
	}
	var result []corev1.Pod
	for _, pod := range pods.Items {
		for _, ownerRef := range pod.GetOwnerReferences() {
			if ownerRef.Kind == "ReplicaSet" && ownerRef.Name == replicaset.GetName() {
				result = append(result, pod)
				break
			}
		}
	}
	return result, nil
}
