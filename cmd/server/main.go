package main

import (
	"flag"
	"net/http"

	"github.com/vladkonst/metrics-alerting/internal/models"
	"github.com/vladkonst/metrics-alerting/routers"
)

func main() {
	addr := &models.NetAddress{Host: "localhost", Port: 8080}
	flag.Var(addr, "a", "Server net address host:port")
	flag.Parse()
	http.ListenAndServe(addr.String(), routers.GetRouter())
}
