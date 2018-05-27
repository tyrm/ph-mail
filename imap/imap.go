package imap

import (
	"log"

	"../util"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func GetClient(address string, username string, password string)  *client.Client {
	imapClient, err := client.DialTLS(address, nil)
	util.PanicOnError(err, "Could not connect to imap server")

	err = imapClient.Login(username, password)
	util.PanicOnError(err, "Could not login to imap server")

	log.Printf("Connected to IMAP server [%s] as [%s]", address, username)
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
		log.Printf("Error geting mailboxes: %v", err)
	}

	return
}

func GetMailbox(imapClient *client.Client, mailboxName string) (mailbox *imap.MailboxStatus, err error) {
	mailbox, err = imapClient.Select(mailboxName, false)
	if err != nil {
		log.Fatal(err)
	}

	return
}
