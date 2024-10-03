package main

import (
	"net/http"

	"github.com/vladkonst/metrics-alerting/routers"
)

func main() {
	http.ListenAndServe("localhost:8080", routers.GetRouter())
}
