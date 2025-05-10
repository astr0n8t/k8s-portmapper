package netstat

import (
	"fmt"
	"testing"
)

func TestNewNetstat(t *testing.T) {
	newListing := NewListing()
	newListing[0] = Listener{}

	for _, l := range GetListening() {
		fmt.Printf("Name: %v Program: %v Port: %v UDP: %v TCP: %v\n", l.Name, l.Program, l.Port, l.UDP, l.TCP)
	}
}
