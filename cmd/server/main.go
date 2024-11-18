package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vladkonst/metrics-alerting/app"
	"github.com/vladkonst/metrics-alerting/internal/configs"
)

func main() {
	done := make(chan bool)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		done <- true
	}()

	cfg := configs.GetServerConfig()
	app := app.NewApp(&done, cfg)
	app.Run()
}
