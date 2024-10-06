package main

import (
	"net/http"

	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/routers"
)

func main() {
	addr := configs.GetServerConfig()
	http.ListenAndServe(addr.String(), routers.GetRouter())
}
