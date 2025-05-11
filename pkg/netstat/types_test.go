package netstat

import (
	"fmt"
	"testing"
)

func TestNewNetstat(t *testing.T) {
	newListing := NewListing()
	newListing[0] = Listener{}

	listing, err := GetListening()
	if err != nil {
		t.Errorf("could not get connections: %v\n", err)
	}

	for _, l := range listing {
		fmt.Printf("Name: %v Program: %v Port: %v UDP: %v TCP: %v\n", l.Name, l.Program, l.Port, l.UDP, l.TCP)
	}
}

func TestNetstatFilter(t *testing.T) {
	newListing := NewListing()
	newListing[0] = Listener{
		Program: "test1",
	}
	newListing[1] = Listener{
		Program: "test2",
	}
	newListing[2] = Listener{
		Program: "nothing",
	}

	newListing.FilterListing("test", "contains")
	if len(newListing) != 2 {
		t.Errorf("filtering list gave incorrect value: %v", newListing)
	}

	newListing.FilterListing("test2", "exact")
	if len(newListing) != 1 {
		t.Errorf("filtering list gave incorrect value: %v", newListing)
	}
}
