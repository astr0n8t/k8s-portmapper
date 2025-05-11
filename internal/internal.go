package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/astr0n8t/k8s-portmapper/pkg/k8s"
	"github.com/astr0n8t/k8s-portmapper/pkg/netstat"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Runs k8s-portmapper
func Run() {
	portMapper := PortMapper{
		Config:  Config(),
		DevMode: true,
	}
	log.Printf("Loaded config file %v", portMapper.Config.ConfigFileUsed())

	log.Println("Loading k8s config")
	if portMapper.Config.GetString("mode") == "production" {
		portMapper.DevMode = false
		portMapper.loadInternalk8sConfig()
		log.Println("Loaded internal k8s config")
	} else {
		portMapper.loadExternalk8sConfig()
		log.Println("Loaded external k8s config")
	}

	// Start main event loop
	go portMapper.startLoop()

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func (p *PortMapper) DebugLog(message string) {
	if p.DevMode {
		log.Printf(message)
	}
}

func (p *PortMapper) startLoop() {
	log.Println("Starting initial checks")
	serviceName := p.Config.GetString("service_name")
	namespace := p.Config.GetString("service_namespace")
	programName := p.Config.GetString("program_name")
	programNameMatch := p.Config.GetString("program_name_match")
	interval := p.Config.GetInt("interval")

	if serviceName == "" || namespace == "" {
		log.Fatalf("service_name and service_namespace are required to be set")
	}

	if programNameMatch != "disabled" && programName == "" {
		log.Fatalf("program_name is required when program_name_match is not disabled")
	}

	if programNameMatch != "disabled" && programNameMatch != "exact" && programNameMatch != "contains" {
		log.Fatalf("program_name_match must be one of disabled, exact, or contains")
	}

	log.Println("Starting k8s-portmapper")
	for {
		p.DebugLog("Getting listening ports from OS")
		listeningPorts, getNetstatErr := netstat.GetListening()
		if getNetstatErr != nil {
			log.Fatalf("could not get netstat: %v\n", getNetstatErr)
		}
		p.DebugLog("Getting service state from k8s")
		updateStateErr := p.updateState(namespace, serviceName)
		if updateStateErr != nil {
			log.Fatalf("could not get service state: %v\n", updateStateErr)
		}
		if programNameMatch != "disabled" {
			p.DebugLog("Filtering netstat listening ports")
			listeningPorts.FilterListing(programName, programNameMatch)
		}
		p.DebugLog("Generating new k8s service port list")
		newPorts := p.generateNewServiceState(&listeningPorts)

		p.DebugLog("Checking to see if port list changed")
		changed := !reflect.DeepEqual(*newPorts, *p.ServiceState)

		if changed {
			log.Printf("%v\n", newPorts)
			p.DebugLog("Port list changed, trying to update k8s service")
			updateServiceErr := p.k8s.SetServicePorts(namespace, serviceName, newPorts)
			if updateServiceErr != nil {
				log.Fatalf("could not update service: %v\n", updateServiceErr)
			}
			p.DebugLog("Updated k8s service")
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (p *PortMapper) updateState(namespace, serviceName string) error {
	var getServiceStateErr error
	p.ServiceState, getServiceStateErr = p.k8s.GetServicePorts(namespace, serviceName)
	if getServiceStateErr != nil {
		return fmt.Errorf("could not get service state: %v\n", getServiceStateErr)
	}

	n := make([]corev1.ServicePort, 0)
	p.NonTrackedPorts = &n

	for _, c := range *p.ServiceState {
		if !strings.HasPrefix(c.Name, "portmap-") {
			p.DebugLog("Adding named port to non tracked state: " + c.Name)
			n := append(*p.NonTrackedPorts, c)
			p.NonTrackedPorts = &n
		}
	}

	return nil
}

func (p *PortMapper) generateNewServiceState(currentListing *netstat.Listing) *[]corev1.ServicePort {
	newListing := netstat.ServicePortSliceToListing(p.NonTrackedPorts)

	var listener netstat.Listener
	var exists bool
	for port, entry := range *currentListing {
		listener, exists = newListing[port]
		if exists {
			// one of them has to be true for it to exist already
			// this shouldn't happen and these ports should be mapped manually
			// or else the port name will be a duplicate and cause errors
			if !listener.TCP && entry.TCP {
				log.Printf("filter allowed port %v on TCP which is mapped manually on UDP but not TCP; will not map this port to prevent errors.  Please map it manually if required")
			} else if !listener.UDP && entry.UDP {
				log.Printf("filter allowed port %v on UDP which is mapped manually on TCP but not UDP; will not map this port to prevent errors.  Please map it manually if required")
			}
			delete(*currentListing, port)
		}
	}

	newState := append(*p.NonTrackedPorts, *currentListing.ListingToServicePortSlice()...)

	return &newState
}

func (p *PortMapper) loadInternalk8sConfig() {
	k8sConfig, k8sErr := k8s.NewK8S()
	if k8sErr != nil {
		log.Fatalf("could not load k8s config: %v", k8sErr)
	}

	p.k8s = *k8sConfig
}

func (p *PortMapper) loadExternalk8sConfig() {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to find k8s config\n")
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create k8s clientset\n")
	}

	k := k8s.NewK8SFrom(*config, *clientset)
	p.k8s = *k
}
