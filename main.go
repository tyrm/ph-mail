package main

import (
	"log"

	"./env"
	"./imap"
	"./models"
	"./util"
)

func main() {
	var myEnv env.Env
	myEnv.Config = env.CollectConfig()

	// Connect IMAP
	myEnv.IMAP = imap.GetClient(myEnv.Config.MailIMAPServer, myEnv.Config.MailUsername, myEnv.Config.MailPassword)
	defer myEnv.IMAP.Logout()

	// Connect DB
	myEnv.DB = models.GetClient(myEnv.Config.DBEngine)
	defer myEnv.DB.Close()

	// Get All Mail Mailbox
	mbox, err := imap.GetMailbox(myEnv.IMAP, "[Gmail]/All Mail")
	util.PanicOnError(err, "Could not login to imap server")

	var cursor int64 = int64(mbox.Messages)
	log.Printf("Messages: %d (%d)", cursor, mbox.Messages)
	for i := cursor; i > 0; i = i - 100 {
		from := uint32(1)
		if i > 99 {from = uint32(i - 99)}

		log.Printf("Range: %d-%d (%d)", i, from, 1 + i - int64(from))
		envelopes, _ := imap.GetEnvelopes(myEnv.IMAP, mbox, from, uint32(i))

		log.Printf("Last %d messages:", 1 + i - int64(from))
		for _, msg := range envelopes {
			models.GetOrCreateEnvelope(myEnv.DB, msg)
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

	log.Println("Done!")
}

