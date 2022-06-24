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

			if ip.To4() != nil && strings.Count(ip.To4().String(), ":") < 2 && !strings.Contains(i.Name, "WSL") && !ip.IsLoopback() {
				ipv4 = ip
			}
		}
	}

	fmt.Println("IPv4 address : ", ipv4.String())
	return ipv4.String()
}

func StartScan(ipv4Arr []string, port int, results *[]ScanResult) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i < 255; i++ {
		wg.Add(1)
		target := fmt.Sprintf("%s.%d", strings.Join(ipv4Arr, "."), i)

		go func(target string) {
			defer wg.Done()
			fmt.Println("Port Scanning " + target)
			open := ScanPort("tcp", target, port)

			*results = append(*results, ScanResult{port, open, target})

		}(target)

	}

}

func ScanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, time.Second)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
