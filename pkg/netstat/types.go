package netstat

type Listing map[uint32]Listener

type Listener struct {
	Name    string
	Program string
	Port    uint32
	TCP     bool
	UDP     bool
}
