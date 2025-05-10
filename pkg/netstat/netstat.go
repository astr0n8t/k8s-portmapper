package netstat

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	corev1 "k8s.io/api/core/v1"
)

const (
	TCP_TYPE = syscall.SOCK_STREAM
	UDP_TYPE = syscall.SOCK_DGRAM
)

func NewListing() Listing {
	var listing Listing
	listing = make(map[uint32]Listener)
	return listing
}

func GetListening() Listing {
	listing := NewListing()

	// Retrieve all connections
	conns, err := net.Connections("all")
	if err != nil {
		fmt.Printf("Error retrieving connections: %v\n", err)
		return nil
	}

	// Process connections
	for _, c := range conns {
		// only care about listening connections
		if (c.Type == TCP_TYPE && c.Status != "LISTEN") ||
			(c.Type == UDP_TYPE && (c.Raddr.IP != "" || c.Raddr.Port != 0)) {
			continue
		}

		var listener Listener
		var exists bool
		listener, exists = listing[c.Laddr.Port]

		if !exists {
			listener = Listener{
				Name:    "",
				Program: "",
				Port:    c.Laddr.Port,
				TCP:     false,
				UDP:     false,
			}
		}

		listener.TCP = (listener.TCP || c.Type == TCP_TYPE)
		listener.UDP = (listener.UDP || c.Type == UDP_TYPE)

		if listener.Program == "" {
			if c.Pid != 0 {
				if p, err := process.NewProcess(c.Pid); err == nil {
					if name, err := p.Name(); err == nil {
						listener.Program = name
					}
				}
			}
		}

		listener.Name = "portmap-" + strconv.Itoa(int(listener.Port))

		listing[c.Laddr.Port] = listener
	}

	return listing
}

func (l *Listing) ListingToServicePortSlice() *[]corev1.ServicePort {
	// start with initial size of map though it could double potentially
	servicePorts := make([]corev1.ServicePort, len(*l))

	for _, p := range *l {
		portInt32 := int32(p.Port)

		servicePort := corev1.ServicePort{
			Name:     p.Name,
			Port:     portInt32,
			Protocol: corev1.ProtocolTCP,
		}

		if p.TCP {
			servicePorts = append(servicePorts, servicePort)
		}

		if p.UDP {
			servicePort.Protocol = corev1.ProtocolUDP
			servicePorts = append(servicePorts, servicePort)
		}
	}

	return &servicePorts
}

func ServiePortSliceToListing(ports *[]corev1.ServicePort) Listing {
	listing := NewListing()

	for _, p := range *ports {
		portUint32 := uint32(p.Port)

		var listener Listener
		var exists bool
		listener, exists = listing[portUint32]

		if !exists {
			listener = Listener{
				Name:    p.Name,
				Program: "",
				Port:    portUint32,
				TCP:     false,
				UDP:     false,
			}
		}

		listener.TCP = (listener.TCP || p.Protocol == corev1.ProtocolTCP)
		listener.UDP = (listener.UDP || p.Protocol == corev1.ProtocolUDP)

		listing[portUint32] = listener
	}

	return listing
}
