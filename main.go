package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/j-keck/arping"
)

type Service struct {
	IsOnline bool
	Count    int
	Hostname string
	Address  string
	HWAdress net.HardwareAddr
}

func main() {
	// Get own IPv4 address, split it into an array and remove the last element.
	ipv4 := GetOwnIPv4Adress()
	ipv4Arr := strings.Split(ipv4, ".")
	ipv4Arr = ipv4Arr[:len(ipv4Arr)-1]

	port := 80

	var results []ScanResult
	services := []Service{}
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		return
	}
	//Scan the IP Range of the local network for open ports.
	StartScan(ipv4Arr, port, &results)

	//Loop over the results and create a Service struct for each open port.
	for _, result := range results {
		if result.State {
			log.Println(result.Address + ":" + strconv.Itoa(result.Port) + " is open")

			dstIP := net.ParseIP(result.Address)
			hwAddr, duration, err := arping.Ping(dstIP)

			if err != nil {
				log.Println(err)
				continue
			}

			services = append(services, Service{IsOnline: true, Count: 0, Hostname: hostname, Address: fmt.Sprintf("http://%s:%d", result.Address, result.Port), HWAdress: hwAddr})

			log.Printf("%s (%s) %d usec\n", dstIP, hwAddr, duration/1000)

		}
	}

	wg := sync.WaitGroup{}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	//Run indefinitely, checking the status of the services.
	for {
		for i := 0; i < len(services); i++ {
			wg.Add(1)
			go getResponse(services[i].Address, &services[i].IsOnline, &services[i].Count, hostname, &wg, &client)
		}
		wg.Wait()

		<-time.After(15 * time.Second)
	}
}

//getResponse is a function that checks the status of a service and sends a notification if the service is down after retrying 3 times.
func getResponse(address string, isOnline *bool, count *int, hostname string, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	_, err := client.Get(address)
	if err != nil {
		if *isOnline {
			*count++
		}
		log.Println("Request failed the: " + strconv.Itoa(*count) + " time,	 for: " + address)
		if *count == 3 && *isOnline {
			log.Println("Request to " + address + " failed 3 times, marking as offline")
			SendMail(address, hostname)
			*isOnline = false
		}

		return
	}
	*isOnline = true
	*count = 0
	log.Println("Successfully pinged " + address)

}
