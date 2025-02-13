package kubecli

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetStatefulSet get the statefulSet with a specific name.
func (k *KubeCli) GetStatefulSet(name string) (*appsv1.StatefulSet, error) {
	statefulSet, err := k.ClientSet.AppsV1().StatefulSets(k.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get StatefulSet: %w", err)
	}
	return statefulSet, nil
}
