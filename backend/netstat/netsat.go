package netstat

import (
	"fmt"

	"github.com/cakturk/go-netstat/netstat"
)

type NetstatResult struct {
	Sock   netstat.SockTabEntry
	IpAddr string
}

var NetstatResults []NetstatResult

func GetSocks() (*[]NetstatResult, error) {
	var err error
	var socks []netstat.SockTabEntry
	var results []NetstatResult = make([]NetstatResult, 0)

	// UDP sockets
	socks, err = netstat.UDPSocks(netstat.NoopFilter)
	if err != nil {
		return nil, err
	}

	for _, e := range socks {
		var entry NetstatResult
		fmt.Printf("%v\n", e)

		if e.RemoteAddr.IP.To4() != nil {
			entry.IpAddr = e.RemoteAddr.IP.String()
			entry.Sock = e
			results = append(results, entry)
		}
	}

	// TCP sockets
	socks, err = netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		return nil, err
	}

	for _, e := range socks {
		var entry NetstatResult
		fmt.Printf("%v\n", e)

		if e.RemoteAddr.IP.To4() != nil {
			entry.IpAddr = e.RemoteAddr.IP.String()
			entry.Sock = e
			results = append(results, entry)
		}
	}

	return &results, nil
}
