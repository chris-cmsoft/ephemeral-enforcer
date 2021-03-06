package main

import (
	"github.com/robfig/cron/v3"
	"github.com/spazzy757/ephemeral-enforcer/pkg/helpers"
	"github.com/spazzy757/ephemeral-enforcer/pkg/workloadkiller"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	kubeconfig := helpers.GetConfig()
	c := cron.New(
		cron.WithLogger(
			cron.VerbosePrintfLogger(
				log.New(os.Stdout, "Ephemeral Enforcer: ", log.LstdFlags),
			),
		),
	)
	_, err := c.AddFunc(os.Getenv("ENFORCER_SCHEDULE"), func() {
		workloadkiller.KillWorkloads(kubeconfig)
	})
	if err != nil {
		log.Fatal("Error", err)
	}
	c.Start()
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	log.Println("Ending Program")
	c.Stop()
}
