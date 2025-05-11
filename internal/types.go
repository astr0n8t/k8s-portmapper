package internal

import (
	"github.com/astr0n8t/k8s-portmapper/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
)

type PortMapper struct {
	Config          ConfigStore
	k8s             k8s.K8S
	NonTrackedPorts *[]corev1.ServicePort
	ServiceState    *[]corev1.ServicePort
	DevMode         bool
}
