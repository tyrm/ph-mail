package imap

import (
	"github.com/juju/loggo"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var logger *loggo.Logger

func GetClient(address string, username string, password string)  *client.Client {
	newLogger :=  loggo.GetLogger("mail.imap")
	logger = &newLogger

	imapClient, err := client.DialTLS(address, nil)
	if err != nil {
		logger.Criticalf("Could not connect to imap server: %s", err)
		panic("PANIC!")
	}

	err = imapClient.Login(username, password)
	if err != nil {
		logger.Criticalf("Could not login to imap server: %s", err)
		panic("PANIC!")
	}

	logger.Infof("Connected to IMAP server [%s] as [%s]", address, username)
	return imapClient
}

func GetEnvelopes(imapClient *client.Client, mbox *imap.MailboxStatus, from uint32, to uint32) (envelopes []*imap.Envelope, err error) {
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messageChannel := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messageChannel)
	}()

	for msg := range messageChannel {
		envelopes = append(envelopes, msg.Envelope)
	}

	err = <-done
	return
}

func GetMailboxes(imapClient *client.Client) (mailboxes []*imap.MailboxInfo, err error) {
	// List mailboxes
	mailboxChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func () {
		done <- imapClient.List("", "*", mailboxChan)
	}()

	for m := range mailboxChan {
		mailboxes = append(mailboxes, m)
	}

	if err := <-done; err != nil {
		logger.Errorf("Error geting mailboxes: %v", err)
	}

	return
}

func GetMailbox(imapClient *client.Client, mailboxName string) (mailbox *imap.MailboxStatus, err error) {
	mailbox, err = imapClient.Select(mailboxName, false)
	if err != nil {
		logger.Errorf("Error getting mailbox: %s", err)
	}

	return
}
