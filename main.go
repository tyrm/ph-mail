package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"./imap"
	"./models"
	"./sync"
	"./web"

	"github.com/juju/loggo"
)

var logger *loggo.Logger

func main() {
	loggo.ConfigureLoggers("<root>=DEBUG")
	newLogger :=  loggo.GetLogger("mail")
	logger = &newLogger

	config := CollectConfig()

	// Connect IMAP
	imap.InitIMAP(config.MailIMAPServer, config.MailUsername, config.MailPassword)
	defer imap.CloseIMAP()

	// Connect DB
	models.InitDB(config.DBEngine, config.ESHost)
	defer models.CloseDB()

	go web.StartWebServer()
	go sync.StartSyncer()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-nch)

	logger.Infof("Done!")
}

