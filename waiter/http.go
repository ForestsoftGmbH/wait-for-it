package waiter

import (
	"log"
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
	resp, err := http.Get(getUrl(w.Waiter))
	if err != nil {
		log.Fatal(err)
		return false
	}
	if (resp.StatusCode != w.Waiter.Status) && (w.Waiter.Status != 0) {
		return false
	}
	return true
}

func getUrl(w *Waiter) string {
	var url string
	if (w.Port == 443) && w.Host == "localhost" {
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
