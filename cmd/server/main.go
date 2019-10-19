package main

import (
	"os"
	"os/signal"
	"time"

	_ "github.com/emmercm/photos/internal/pkg/logger"
	"github.com/emmercm/photos/internal/pkg/monitor"
	"github.com/emmercm/photos/internal/pkg/store/sqlite3"
	"github.com/emmercm/photos/internal/pkg/transport/http"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := sqlite3.Migrate(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	monitorCollection := monitor.NewCollection()
	monitorCollection.WatchDirectory("/Users/christianemmer/Resilio Sync/DCIM")
	monitorCollection.UnwatchDirectory("/Users/christianemmer/Resilio Sync/DCIM")
	monitorCollection.WatchDirectory("/Users/christianemmer/Resilio Sync/DCIM")
	monitorCollection.WatchDirectory("./")

	go monitorCollection.Start()
	go http.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		os.Exit(0)
	}()

	for {
		time.Sleep(time.Second)
	}
}
