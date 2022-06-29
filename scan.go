package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port    int
	State   bool
	Address string
}

// GetOwnIPv4Adress returns the IPv4 address of the local machine.
func GetOwnIPv4Adress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var ipv4 net.IP

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		for index, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP

			case *net.IPAddr:
				ip = v.IP
			}
			log.Println("[", index, "]", i.Name, ">", ip)
			if ip.To4() != nil &&
				strings.Count(ip.To4().String(), ":") < 2 &&
				!strings.Contains(i.Name, "WSL") &&
				!ip.IsLoopback() {
				ipv4 = ip
			}
		}
	}

	log.Println("IPv4 address : ", ipv4.String())
	return ipv4.String()
}

//StartScan starts a scan of the local network for open ports.
func StartScan(ipv4Arr []string, port int, results *[]ScanResult) {
	log.Println("Scanning network for open ports...")
	log.Println("IP Range : ", strings.Join(append(ipv4Arr, "1"), "."), "---", strings.Join(append(ipv4Arr, "255"), "."))
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i < 255; i++ {

		wg.Add(1)
		target := fmt.Sprintf("%s.%d", strings.Join(ipv4Arr, "."), i)

		go func(target string) {
			defer wg.Done()
			open := ScanPort("tcp", target, port)

			*results = append(*results, ScanResult{port, open, target})

		}(target)

	}

}

//ScanPort scans a port on a target host.
func ScanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, time.Second)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
