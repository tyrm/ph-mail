package main

import (
	"./imap"
	"./models"
	"github.com/juju/loggo"
)

var logger *loggo.Logger

func main() {
	loggo.ConfigureLoggers("<root>=DEBUG")
	newLogger :=  loggo.GetLogger("mail")
	logger = &newLogger

	config := CollectConfig()

	// Connect IMAP
	//imapCon := imap.GetClient(config.MailIMAPServer, config.MailUsername, config.MailPassword)
	imap.GetClient(config.MailIMAPServer, config.MailUsername, config.MailPassword)
	//defer imapCon.Logout()

	// Connect DB
	//myEnv.DB = models.GetClient(myEnv.Config.DBEngine)
	models.InitDB(config.DBEngine, config.ESHost)
	defer models.CloseDB()

	// Get All Mail Mailbox
	mbox, err := imap.GetMailbox("[Gmail]/All Mail")
	if err != nil {
		logger.Criticalf("Could not login to imap server: %s", err)
		panic("PANIC!")
	}

	var cursor int64 = int64(mbox.Messages)
	logger.Infof("Messages: %d (%d)", cursor, mbox.Messages)
	//log.Printf("Messages: %d (%d)", cursor, mbox.Messages)
	for i := cursor; i > 0; i = i - 100 {
		from := uint32(1)
		if i > 99 {from = uint32(i - 99)}

		logger.Debugf("Range: %d-%d (%d)", i, from, 1 + i - int64(from))
		envelopes, _ := imap.GetEnvelopes(mbox, from, uint32(i))

		logger.Debugf("Last %d messages:", 1 + i - int64(from))
		for _, msg := range envelopes {
			models.GetOrCreateEnvelope(msg)
			//log.Printf(" Envelope:\n [%T] [%s]", envelope, envelope)
		}
	}
















	/*
	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 99 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 99
	}
	envelopes, err := imap.GetEnvelopes(env.imap, mbox, from, to)

	log.Println("Last 4 messages:")
	for _, msg := range envelopes {
		log.Printf("* %s", reflect.TypeOf(msg.To))
		log.Printf("* %s", msg)
		envelope, _ := db.GetOrCreateEnvelope(env.db, msg)
		log.Printf("Address:\n [%s] [%s]", envelope, )
	}*/

	logger.Infof("Done!")
}

