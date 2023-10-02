package waiter

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
)

type HttpWaiter struct {
	Waiter *Waiter
}

func NewHttpWaiter(w *Waiter) HttpWaiter {
	return HttpWaiter{
		Waiter: w,
	}
}

func (w HttpWaiter) ShouldExecute() bool {
	switch w.Waiter.Port {
	case 80, 443, 8080, 8085:
		return true

	}
	return false
}

func (w HttpWaiter) IsReady() bool {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(getUrl(w.Waiter))
	if err != nil {
		fmt.Println(err)
		return false
	}
	if (resp.StatusCode != w.Waiter.Status) && (w.Waiter.Status != 0) {
		fmt.Println("Status code is " + strconv.Itoa(resp.StatusCode) + " but expected " + strconv.Itoa(w.Waiter.Status))
		return false
	}
	return true
}

func getUrl(w *Waiter) string {
	var url string
	if (w.Port == 443) || (w.Port == 8443) {
		url = "https://"
	} else {
		url = "http://"
	}
	url += w.Host
	if (w.Port != 80) && (w.Port != 443) {
		url += ":" + strconv.Itoa(w.Port)
	}
	url += w.Path
	return url
}
