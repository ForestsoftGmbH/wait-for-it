package main

import (
	"errors"
	"flag"
	"github.com/ForestsoftGmbH/wait-for-it/waiter"
	"os"
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
	flag.DurationVar(&timeout, "timeout", 30, "Timeout in seconds for the wait")

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
	done := make(chan bool)
	var result bool
	var err error

	to := time.NewTimer(timeout)
	defer to.Stop()
	defer ticker.Stop()
	httpWaiter := waiter.NewHttpWaiter(w)
	if httpWaiter.ShouldExecute() {
		go func(waiter waiter.HttpWaiter, timer *time.Timer) {
			for {
				select {
				case <-timer.C:
					err = errors.New("Timeout reached")
					result = false
				case <-done:
					result = true
					err = nil
				case <-ticker.C:
					if waiter.IsReady() {
						done <- true
					}
				}
			}
		}(httpWaiter, to)
	}
	return result, err
}
