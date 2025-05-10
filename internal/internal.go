package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/astr0n8t/k8s-portmapper/pkg/netstat"
)

// Runs k8s-portmapper
func Run() {
	// Make sure we can load config
	config := Config()
	log.Printf("Loaded config file %v", config.ConfigFileUsed())

	newListing := netstat.GetListening()

	for _, l := range newListing {
		fmt.Printf("Program: %v Port: %v UDP: %v TCP: %v\n", l.Program, l.Port, l.UDP, l.TCP)
	}

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop
}
