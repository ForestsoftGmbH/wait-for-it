package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ForestsoftGmbH/wait-for-it/waiter"
	"os"
	"sync"
	"time"
)

func main() {
	var port, status int
	var timeout time.Duration
	var host, path string

	flag.IntVar(&port, "p", 80, "Provide a port number")
	flag.StringVar(&host, "host", "localhost", "Hostname to check")
	flag.StringVar(&path, "path", "/", "Path to check")
	flag.IntVar(&status, "statusCode", 200, "Check for status code")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Timeout in seconds for the wait")

	flag.Parse()

	result, err := WaitForIt(host, path, port, status, timeout)
	if err != nil {
		panic(err)
	}
	if result {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func WaitForIt(host, path string, port, status int, timeout time.Duration) (bool, error) {
	w := waiter.NewWaiter(host, path, port, status)

	ticker := time.NewTicker(time.Second)

	var result bool
	var err error

	to := time.NewTimer(timeout)
	defer to.Stop()
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	httpWaiter := waiter.NewHttpWaiter(w)
	if httpWaiter.ShouldExecute() {
		go func(waiter waiter.HttpWaiter, timer *time.Timer, timeout time.Duration) {
			port := fmt.Sprintf("%d", waiter.Waiter.Port)
			fmt.Println("Waiting for " + waiter.Waiter.Host + ":" + port + waiter.Waiter.Path + " to be ready")
			defer wg.Done()
			for {
				select {
				case <-timer.C:
					err = errors.New(fmt.Sprintf("Timeout %s reached", timeout))
					result = false
					return
				case <-ticker.C:
					if waiter.IsReady() {
						result = true
						err = nil
						fmt.Println("Service is ready")
						return
					}
				}
			}

		}(httpWaiter, to, timeout)
	}
	wg.Wait()
	return result, err
}
