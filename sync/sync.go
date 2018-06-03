package sync

import (
	"time"

	"../imap"
	"../models"
	"github.com/juju/loggo"
)

var logger *loggo.Logger
var mailboxName = "[Gmail]/All Mail"
var syncing = false

func StartSyncer() {
	newLogger :=  loggo.GetLogger("mail.sync")
	logger = &newLogger

	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			mbox, err := imap.GetMailbox(mailboxName)
			if err != nil {
				logger.Errorf("Could not get folder %s: %s", "[Gmail]/All Mail", err)
				return
			}

			LastEnvelope, err := imap.GetLastEnvelope(mbox)
			if err != nil {
				logger.Errorf("Could not get envelope: %s", err)
				return
			}

			if models.EnvelopeExistsByMsgID(LastEnvelope.MessageId) {
				logger.Debugf("Envelope %s seen before.", LastEnvelope.MessageId)
			} else {
				logger.Debugf("New envelope %s. Starting Sync.", LastEnvelope.MessageId)
				SyncRecentEnvelopes()
			}
		}
	}()

}

func SyncAllEnvelopes() {
	if !syncing {
		syncing = true

		// Get All Mail Mailbox
		mbox, err := imap.GetMailbox(mailboxName)
		if err != nil {
			logger.Criticalf("Could not login to imap server: %s", err)
			panic("PANIC!")
		}

		cursor := int64(mbox.Messages)
		logger.Infof("Messages: %d (%d)", cursor, mbox.Messages)

		for i := cursor; i > 0; i = i - 100 {
			from := uint32(1)
			if i > 99 {from = uint32(i - 99)}

			logger.Debugf("Range: %d-%d (%d)", i, from, 1 + i - int64(from))
			envelopes, _ := imap.GetEnvelopes(mbox, from, uint32(i))

			logger.Debugf("Last %d messages:", 1 + i - int64(from))
			for _, msg := range envelopes {
				if !models.EnvelopeExistsByMsgID(msg.MessageId) {
					models.CreateEnvelope(msg)
				}
			}
		}

		syncing = false
	} else {
		logger.Infof("Ignoring sync request because sync is already in progress.")
	}
}

func SyncRecentEnvelopes() {
	if !syncing {
		syncing = true

		// Get All Mail Mailbox
		mbox, err := imap.GetMailbox(mailboxName)
		if err != nil {
			logger.Criticalf("Could not login to imap server: %s", err)
			panic("PANIC!")
		}

		cursor := int64(mbox.Messages)
		logger.Infof("Messages: %d (%d)", cursor, mbox.Messages)

		for i := cursor; i > 0; i = i - 20 {
			from := uint32(1)
			if i > 19 {from = uint32(i - 19)}

			logger.Debugf("Range: %d-%d (%d)", i, from, 1 + i - int64(from))
			envelopes, _ := imap.GetEnvelopes(mbox, from, uint32(i))

			logger.Debugf("Last %d messages:", 1 + i - int64(from))

			var foundNew = false
			for _, msg := range envelopes {
				if !models.EnvelopeExistsByMsgID(msg.MessageId) {
					foundNew = true
					models.CreateEnvelope(msg)
				}
			}

			// If we didn't find a new message this round break.
			if !foundNew {break}
		}

		syncing = false
	} else {
		logger.Infof("Ignoring sync request because sync is already in progress.")
	}
}