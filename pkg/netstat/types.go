package netstat

type Netstat struct {
	Listing map[uint32]Listener
}

type Listener struct {
	Program string
	Port    uint32
	TCP     bool
	UDP     bool
}
