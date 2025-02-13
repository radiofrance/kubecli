package kubecli

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeployment get the deployment with a specific name.
func (k *KubeCli) GetDeployment(name string) (*appsv1.Deployment, error) {
	deployment, err := k.ClientSet.AppsV1().Deployments(k.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Deployment: %w", err)
	}
	return deployment, nil
}
