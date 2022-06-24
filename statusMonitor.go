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
)

//struct that matches getResponse input
type Service struct {
	IsOnline bool
	Count    int
	Hostname string
	Address  string
}

func main() {
	ipv4 := GetOwnIPv4Adress()
	ipv4Arr := strings.Split(ipv4, ".")

	ipv4Arr = ipv4Arr[:len(ipv4Arr)-1]
	var results []ScanResult
	StartScan(ipv4Arr, 80, &results)
	for _, result := range results {
		if result.State {
			fmt.Println(result.Address + ":" + strconv.Itoa(result.Port) + " is open")
		}
	}

	hostname, err := os.Hostname()
	var service1 = Service{
		IsOnline: true,
		Count:    0,
		Hostname: hostname,
		Address:  fmt.Sprintf("http://%s:%d", results[0].Address, results[0].Port),
	}
	var service2 = Service{
		IsOnline: true,
		Count:    0,
		Hostname: hostname,
		Address:  fmt.Sprintf("http://%s:%d", results[1].Address, results[1].Port),
	}

	if err != nil {
		log.Println(err)
		return
	}
	wg := sync.WaitGroup{}
	for {

		go getResponse(service1.Address, &service1.IsOnline, &service1.Count, hostname, wg)
		go getResponse(service2.Address, &service2.IsOnline, &service2.Count, hostname, wg)

		wg.Wait()
		log.Println(service1.IsOnline, service1.IsOnline)
		log.Println(service2.Count, service2.Count)

		<-time.After(5 * time.Second)
	}
}

func getResponse(address string, isOnline *bool, count *int, hostname string, wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	_, err := http.Get(address)
	if err != nil {
		log.Println("Request failed for: " + address)
		if *isOnline {
			*count++
		}
		if *count == 3 && *isOnline {
			log.Println("send mail")
			SendMail(address, hostname)
			*count = 0
			*isOnline = false
		}
		return
	}
	log.Println("Successfully pinged " + address)
}
