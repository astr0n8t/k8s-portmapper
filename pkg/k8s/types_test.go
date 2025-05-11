package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNewK8S(t *testing.T) {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to find config, Will not test k8s\n")
		return
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("failed to create clientset: %v", err)
	}

	k := NewK8SFrom(*config, *clientset)

	ports, portGetErr := k.GetServicePorts("testing", "test-service")
	if portGetErr != nil {
		t.Errorf("failed to get service ports: %v", portGetErr)
	}

	newPorts := append(*ports, corev1.ServicePort{
		Name:     "test",
		Port:     8080,
		Protocol: corev1.ProtocolTCP,
	})

	portSetErr := k.SetServicePorts("testing", "test-service", &newPorts)
	if portSetErr != nil {
		t.Errorf("failed to set service ports: %v", portSetErr)
	}

	portResetErr := k.SetServicePorts("testing", "test-service", ports)
	if portResetErr != nil {
		t.Errorf("failed to reset service ports: %v", portResetErr)
	}
}
