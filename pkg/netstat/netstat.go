package netstat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

func GetListening() (Listing, error) {
	listing := NewListing()

	// Retrieve all connections
	conns, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve connections: %v\n", err)
	}

	// Process connections
	for _, c := range conns {
		// only care about listening connections
		if (c.Type == TCP_TYPE && c.Status != "LISTEN") ||
			(c.Type == UDP_TYPE && (c.Raddr.IP != "" && c.Raddr.Port != 0)) {
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
					} else {
						log.Printf("issue retrieving listener name: %v\n", err)
					}
				}
			}
		}

		listener.Name = "portmap-" + strconv.Itoa(int(listener.Port))

		listing[c.Laddr.Port] = listener
	}

	return listing, nil
}

func (l *Listing) FilterListing(program, method string) {
	for p, e := range *l {
		if method == "exact" && e.Program != program {
			delete(*l, p)
		} else if method == "contains" && !strings.Contains(e.Program, program) {
			delete(*l, p)
		}
	}
}

func (l *Listing) ListingToServicePortSlice() *[]corev1.ServicePort {
	// start with initial size of map though it could double potentially
	servicePorts := make([]corev1.ServicePort, 0)

	for _, p := range *l {
		portInt32 := int32(p.Port)

		tracked := strings.HasPrefix(p.Name, "portmap-")

		servicePort := corev1.ServicePort{
			Name:     p.Name,
			Port:     portInt32,
			Protocol: corev1.ProtocolTCP,
		}

		if p.TCP {
			if tracked {
				servicePort.Name = "portmap-" + strconv.Itoa(int(servicePort.Port)) + "-t"
			}
			servicePorts = append(servicePorts, servicePort)
		}

		if p.UDP {
			if tracked {
				servicePort.Name = "portmap-" + strconv.Itoa(int(servicePort.Port)) + "-u"
			}
			servicePort.Protocol = corev1.ProtocolUDP
			servicePorts = append(servicePorts, servicePort)
		}
	}

	return &servicePorts
}

func ServicePortSliceToListing(ports *[]corev1.ServicePort) Listing {
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
