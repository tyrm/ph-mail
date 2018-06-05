package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/juju/loggo"
)

var imapClient *client.Client
var logger *loggo.Logger

func CloseIMAP() {
	imapClient.Logout()
	return
}

func TraceEnvelope(m *imap.Message) {
	msg := "Envelope Object:\nSeqNum: %v\nItems: %v\nEnvelope: %v\nBodyStructure: %v\nFlags: %v\nInternalDate: %v\nSize: %v\nUid: %v\nBody: %v"
	logger.Tracef(msg, m.SeqNum, m.Items, m.Envelope, m.BodyStructure, m.Flags, m.InternalDate, m.Size, m.Uid, m.Body)
}

func GetBodies(msg *imap.Message) {

}

func GetMessages(mbox *imap.MailboxStatus, from uint32, to uint32) (messages []*imap.Message, err error) {
	client := imapClient

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messageChannel := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchFlags}, messageChannel)
	}()

	for msg := range messageChannel {
		messages = append(messages, msg)
	}

	err = <-done
	return
}

func GetLastMessage(mbox *imap.MailboxStatus) (message *imap.Message, err error) {
	LastMsgNumber := mbox.Messages
	LastEnvelope, err := GetMessages(mbox, LastMsgNumber, LastMsgNumber)
	if err == nil {
		message = LastEnvelope[0]
	}

	return
}

func GetMailbox(mailboxName string) (mailbox *imap.MailboxStatus, err error) {
	client := imapClient
	mailbox, err = client.Select(mailboxName, true)
	if err != nil {
		logger.Errorf("Error getting mailbox: %s", err)
	}

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

func InitIMAP(address string, username string, password string) {
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
