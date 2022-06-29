package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Service struct {
	IsOnline     bool
	Count        int
	Hostname     string
	Address      string
	OfflineSince time.Time
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
		return
	}

	//Get own IPv4 address, split it into an array and remove the last element.
	//This is just a safety incase the IP address range is not the standard 192.168.1.0/24
	ipv4 := GetOwnIPv4Adress()
	ipv4Arr := strings.Split(ipv4, ".")
	ipv4Arr = ipv4Arr[:len(ipv4Arr)-1]

	port := 3000

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

			if err != nil {
				log.Println(err)
				continue
			}

			services = append(services, Service{IsOnline: true, Count: 0, Hostname: hostname, Address: fmt.Sprintf("http://%s:%d", result.Address, result.Port)})

		}
	}

	wg := sync.WaitGroup{}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	if len(services) == 0 {
		log.Println("No services found")
		return
	}
	//Run indefinitely, checking the status of the services.

	for {
		for i := 0; i < len(services); i++ {
			wg.Add(1)
			go getResponse(&services[i], &wg, &client)
		}
		wg.Wait()

		<-time.After(15 * time.Second)
	}
}

//getResponse is a function that checks the status of a service and sends a notification if the service is down after retrying 3 times.
func getResponse(service *Service, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	_, err := client.Get(service.Address)
	if err != nil {
		if service.IsOnline {
			service.Count++
			log.Println("Request failed the: " + strconv.Itoa(service.Count) + ". time, for: " + service.Address)
		}

		if service.Count == 3 && service.IsOnline {
			log.Println("Request to " + service.Address + " failed 3 times, marking as offline")
			SendMail(service.Address, service.Hostname)
			service.OfflineSince = time.Now()
			service.IsOnline = false
			return
		}

		//log that the service is still offline.
		if !service.IsOnline {
			log.Println("Request to " + service.Address + " failed, service is still offline")
			log.Println("Service has been offline since: ", service.OfflineSince)
		}
		return
	}
	service.IsOnline = true
	service.Count = 0
	log.Println("Successfully pinged " + service.Address)

}
