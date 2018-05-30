package imap

import (
	"github.com/juju/loggo"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var imapClient *client.Client
var logger *loggo.Logger

func GetClient(address string, username string, password string) {
	newLogger :=  loggo.GetLogger("mail.imap")
	logger = &newLogger

	newImapClient, err := client.DialTLS(address, nil)
	if err != nil {
		logger.Criticalf("Could not connect to imap server: %s", err)
		panic("PANIC!")
	}
	imapClient = newImapClient

	err = imapClient.Login(username, password)
	if err != nil {
		logger.Criticalf("Could not login to imap server: %s", err)
		panic("PANIC!")
	}

	logger.Infof("Connected to IMAP server [%s] as [%s]", address, username)
	return
}

func GetEnvelopes(mbox *imap.MailboxStatus, from uint32, to uint32) (envelopes []*imap.Envelope, err error) {
	client := imapClient

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messageChannel := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messageChannel)
	}()

	for msg := range messageChannel {
		envelopes = append(envelopes, msg.Envelope)
	}

	err = <-done
	return
}

func GetMailboxes() (mailboxes []*imap.MailboxInfo, err error) {
	client := imapClient

	// List mailboxes
	mailboxChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func () {
		done <- client.List("", "*", mailboxChan)
	}()

	for m := range mailboxChan {
		mailboxes = append(mailboxes, m)
	}

	if err := <-done; err != nil {
		logger.Errorf("Error geting mailboxes: %v", err)
	}

	return
}

func GetMailbox(mailboxName string) (mailbox *imap.MailboxStatus, err error) {
	client := imapClient
	mailbox, err = client.Select(mailboxName, false)
	if err != nil {
		logger.Errorf("Error getting mailbox: %s", err)
	}

	return
}
