package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8S struct {
	config    rest.Config
	clientset kubernetes.Clientset
}
