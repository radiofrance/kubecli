package kubecli

import (
	"context"
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetReplicaSet get the replicaSet whose owner is a specific deployment.
func (k *KubeCli) GetReplicaSet(deployment *appsv1.Deployment) (*appsv1.ReplicaSet, error) {
	replicaSets, err := k.ClientSet.
		AppsV1().ReplicaSets(deployment.GetNamespace()).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list replicasets: %w", err)
	}

	var mostRecentRS *appsv1.ReplicaSet
	for _, replicaSet := range replicaSets.Items {
		for _, ownerRef := range replicaSet.GetOwnerReferences() {
			if ownerRef.Kind == "Deployment" && ownerRef.Name == deployment.Name {
				rs := replicaSet
				mostRecentRS, err = getMostRecentRS(mostRecentRS, &rs)
				if err != nil {
					return nil, fmt.Errorf("could not retrieve replicatset: %w", err)
				}
			}
		}
	}
	if mostRecentRS == nil {
		return nil, fmt.Errorf("not found: %w", errors.NewNotFound(schema.GroupResource{
			Group:    deployment.GroupVersionKind().Group,
			Resource: deployment.Name,
		}, deployment.Name))
	}
	return mostRecentRS, nil
}

func getReplicaSetRevision(replicaSet *appsv1.ReplicaSet) (int, error) {
	revisionString := replicaSet.Annotations["deployment.kubernetes.io/revision"]
	revision, err := strconv.Atoi(revisionString)
	if err != nil {
		return 0, fmt.Errorf("error parsing replicaset version: %w", err)
	}
	return revision, nil
}

// GetReplicaSetByName get the replicaSet whose owner is a specific deployment.
func (k *KubeCli) GetReplicaSetByName(deploymentName string) (*appsv1.ReplicaSet, error) {
	deployment, err := k.GetDeployment(deploymentName)
	if err != nil {
		return nil, fmt.Errorf("could not get deployment: %w", err)
	}
	return k.GetReplicaSet(deployment)
}

func getMostRecentRS(rs1, rs2 *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	switch {
	case rs1 == nil && rs2 == nil:
		return nil, nil //nolint:nilnil
	case rs1 == nil:
		return rs2, nil
	case rs2 == nil:
		return rs1, nil
	default:
		rs1Rev, err := getReplicaSetRevision(rs1)
		if err != nil {
			return nil, err
		}

		rs2Rev, err := getReplicaSetRevision(rs2)
		if err != nil {
			return nil, err
		}

		if rs1Rev >= rs2Rev {
			return rs1, nil
		}
		return rs2, nil
	}
}
