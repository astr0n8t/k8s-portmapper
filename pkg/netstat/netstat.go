package netstat

import (
	"fmt"
	"syscall"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

const (
	TCP_TYPE = syscall.SOCK_STREAM
	UDP_TYPE = syscall.SOCK_DGRAM
)

func NewListing() *Netstat {
	listing := Netstat{
		Listing: make(map[uint32]Listener),
	}

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
		listener, exists = listing.Listing[c.Laddr.Port]

		if !exists {
			listener = Listener{
				Program: "",
				Port:    0,
				TCP:     false,
				UDP:     false,
			}
		}

		listener.Port = c.Laddr.Port
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

		listing.Listing[c.Laddr.Port] = listener
	}

	return &listing
}
