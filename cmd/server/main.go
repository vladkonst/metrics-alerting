package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vladkonst/metrics-alerting/app"
)

func main() {
	done := make(chan bool)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		done <- true
	}()
	app, err := app.NewApp(&done)
	if err != nil {
		log.Panic(err)
	}
	app.Run()
}
