package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeout = time.Second * 3
	workers = runtime.NumCPU() * 10
)

const (
	portRangeEnd = 65535
)

func getPorts(r string) []string {
	s := strings.Split(r, "-")
	start, err := strconv.Atoi(s[0])
	if err != nil || start > portRangeEnd || start < 1 {
		return nil
	}
	if len(s) == 1 {
		return []string{s[0]}
	}
	if len(s) > 2 {
		return nil
	}
	end, err := strconv.Atoi(s[1])
	if err != nil || start > portRangeEnd || start < 1 {
		return nil
	}
	toReturn := []string{}
	for i := int64(start); i <= int64(end); i++ {
		toReturn = append(toReturn, strconv.FormatInt(i, 10))
	}
	return toReturn
}

func getAddresses() []string {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments. Usage: portcheck HOST [port|port-range|port1,port2,...]")
	}
	host := os.Args[1]
	addresses := []string{}
	if len(os.Args) == 2 {
		for i := range portRangeEnd {
			if i == 0 {
				continue
			}
			addresses = append(
				addresses,
				net.JoinHostPort(host, strconv.FormatInt(int64(i), 10)))
		}
	}
	if len(os.Args) > 2 {
		ports := os.Args[2]
		for i := range strings.SplitSeq(ports, ",") {
			if r := getPorts(i); r != nil {
				addresses = append(addresses, func(r []string) []string {
					toReturn := []string{}
					for _, j := range r {
						toReturn = append(toReturn, net.JoinHostPort(host, j))
					}
					return toReturn
				}(r)...)
			}
		}
	}
	return addresses
}

func main() {
	workerChan := make(chan struct{}, workers)
	addresses := getAddresses()
	wg := sync.WaitGroup{}
	for _, address := range addresses {
		workerChan <- struct{}{}
		wg.Go(func() {
			conn, _ := net.DialTimeout("tcp", address, timeout)
			defer func() { <-workerChan }()
			if conn != nil {
				_, _ = fmt.Fprintf(os.Stdout, "SUCCESS: %s\n", address)
				if errE := conn.Close(); errE != nil {
					fmt.Fprintf(os.Stderr, "error closing connection: %s\n", errE)
				}
			}
		})
	}
	wg.Wait()
}
