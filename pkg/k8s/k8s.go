package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewK8S() (*K8S, error) {
	// Create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %v", err)
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	return &K8S{
		config:    *config,
		clientset: *clientset,
	}, nil
}

func NewK8SFrom(config rest.Config, clientset kubernetes.Clientset) *K8S {
	return &K8S{
		config:    config,
		clientset: clientset,
	}
}

func (k *K8S) GetServicePorts(namespace, serviceName string) (*[]corev1.ServicePort, error) {
	// create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get the current service
	service, err := k.clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	return &service.Spec.Ports, nil
}

func (k *K8S) SetServicePorts(namespace, serviceName string, ports *[]corev1.ServicePort) error {
	// create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create patch data
	patch := struct {
		Spec struct {
			Ports []corev1.ServicePort `json:"ports"`
		} `json:"spec"`
	}{}
	patch.Spec.Ports = *ports

	patchData, marshallErr := json.Marshal(patch)
	if marshallErr != nil {
		return fmt.Errorf("failed to marshal ports: %v", marshallErr)
	}

	// Apply the patch
	_, patchErr := k.clientset.CoreV1().Services(namespace).Patch(
		ctx,
		serviceName,
		types.MergePatchType,
		patchData,
		metav1.PatchOptions{},
	)
	if patchErr != nil {
		return fmt.Errorf("failed to patch service: %v", patchErr)
	}

	return nil
}
