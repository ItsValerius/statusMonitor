package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service struct {
	IsOnline     bool
	Count        int
	Hostname     string
	Address      string
	OfflineSince time.Time
}

var (
	MAIL_PW   string
	MAIL_ADDR string
)

func init() {
	flag.StringVar(&MAIL_PW, "p", "", "Mail Password")
	flag.StringVar(&MAIL_ADDR, "m", "", "Mail Address")
	flag.Parse()
}

func main() {

	startHour := 6
	endHour := 22

	fmt.Println(time.Now().Hour() > startHour && time.Now().Hour() < endHour)

	//Get own IPv4 address, split it into an array and remove the last element.
	//This is just a safety incase the IP address range is not the standard 192.168.1.0/24
	ipv4 := GetOwnIPv4Adress()
	ipv4Arr := strings.Split(ipv4, ".")
	ipv4Arr = ipv4Arr[:len(ipv4Arr)-1]

	port := 3737

	var results []ScanResult
	services := []Service{}
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		return
	}
	//Scan the IP Range of the local network for open ports.
	//try 3 times before exiting.
	try := 0
	for {

		StartScan(ipv4Arr, port, &results)

		//Loop over the results and create a Service struct for each open port.
		for _, result := range results {
			if result.State {
				log.Println(result.Address + ":" + strconv.Itoa(result.Port) + " is open")

				if err != nil {
					log.Println(err)
					continue
				}

				services = append(services, Service{IsOnline: true, Count: 0, Hostname: hostname, Address: fmt.Sprintf("https://%s:%d", result.Address, result.Port)})

			}
		}
		if len(services) > 0 {
			break
		}
		if try >= 3 {
			log.Println("No services found after 3 tries")
			log.Println("Exiting")
			return
		}
		try++
		log.Println("No services found")
		log.Println("trying again in 30 seconds")
		<-time.After(time.Second * 30)
	}
	wg := sync.WaitGroup{}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}

	//Run indefinitely, checking the status of the services.

	for {
		if time.Now().Hour() > startHour && time.Now().Hour() < endHour {
			for i := 0; i < len(services); i++ {
				wg.Add(1)
				go getResponse(&services[i], &wg, &client)
			}
			wg.Wait()

			<-time.After(15 * time.Second)
		}
		<-time.After(time.Hour)
	}
}

// getResponse is a function that checks the status of a service and sends a notification if the service is down after retrying 3 times.
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
			SendMail(service.Address)
			service.OfflineSince = time.Now()
			service.IsOnline = false
			return
		}

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
