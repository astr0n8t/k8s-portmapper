package netstat

import (
	"fmt"
	"testing"
)

func TestSetAndGetIndex(t *testing.T) {
	newListing := NewListing()
	for _, l := range newListing.Listing {
		fmt.Printf("Program: %v Port: %v UDP: %v TCP: %v\n", l.Program, l.Port, l.UDP, l.TCP)
	}
}
