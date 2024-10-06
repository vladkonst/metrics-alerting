package main

import (
	"net/http"

	"github.com/vladkonst/metrics-alerting/internal/flags"
	"github.com/vladkonst/metrics-alerting/routers"
)

func main() {
	addr := flags.GetNetAddress()
	http.ListenAndServe(addr.String(), routers.GetRouter())
}
