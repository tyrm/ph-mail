package main

import (
	"log"
	"reflect"

	"./config"
	"./db"
	"./imap"
	"./util"

	"github.com/emersion/go-imap/client"
	"github.com/jinzhu/gorm"
)

type Env struct {
	config *config.Config
	db     *gorm.DB
	imap   *client.Client
}

func main() {
	var env Env
	myConfig := config.CollectConfig()

	// Connect IMAP
	env.imap = imap.GetClient(myConfig.MailIMAPServer, myConfig.MailUsername, myConfig.MailPassword)
	defer env.imap.Logout()

	// Connect DB
	env.db = db.GetClient(myConfig.DBEngine)
	defer env.db.Close()



















	// Get All Mail Mailbox
	mbox, err := imap.GetMailbox(env.imap, "[Gmail]/All Mail")
	util.PanicOnError(err, "Could not login to imap server")

	log.Println("Flags for INBOX:", mbox.Messages)
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
	}

	log.Println("Done!")
}