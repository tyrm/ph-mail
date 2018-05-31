package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/emersion/go-imap"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// A message.
type Envelope struct {
	gorm.Model

	MessageId string    `gorm:"unique;index;not null"`
	InReplyTo string
	Date      time.Time
	Subject   string

	From      []Address `gorm:"many2many:mail_efroms"`
	Sender    []Address `gorm:"many2many:mail_esenders"`
	ReplyTo   []Address `gorm:"many2many:mail_ereplytos"`
	To        []Address `gorm:"many2many:mail_etos"`
	Cc        []Address `gorm:"many2many:mail_eccs"`
	Bcc       []Address `gorm:"many2many:mail_ebccs"`
}

type esEnvelope struct {
	MessageId string    `json:"message_id"`
	InReplyTo string    `json:"in_reply_to,omitempty"`
	Date      time.Time `json:"received"`
	Subject   string    `json:"subject"`

	From      []string  `json:"from,omitempty"`
	Sender    []string  `json:"sender,omitempty"`
	ReplyTo   []string  `json:"reply_to,omitempty"`
	To        []string  `json:"to,omitempty"`
	Cc        []string  `json:"cc,omitempty"`
	Bcc       []string  `json:"bcc,omitempty"`
}

const ESDocEnvelope = `{
    "settings": {
        "index": {
            "number_of_shards": "1",
            "number_of_replicas": "0"
        }
    },
    "mappings": {
        "doc": {
            "properties": {
                "bcc": {
                    "type": "text"
                },
                "cc": {
                    "type": "text"
                },
                "from": {
                    "type": "text"
                },
                "in_reply_to": {
                    "type": "text"
                },
                "message_id": {
                    "type": "text"
                },
                "received": {
                    "type": "date"
                },
                "reply_to": {
                    "type": "text"
                },
                "sender": {
                    "type": "text"
                },
                "subject": {
                    "type": "text"
                },
                "to": {
                    "type": "text"
                }
            }
        }
    }
}`

func GetOrCreateEnvelope(imapEnvelope *imap.Envelope) (envelope Envelope, err error) {
	dbErr := db.Preload("From").Preload("Sender").Preload("ReplyTo").
		Preload("To").Preload("Cc").Preload("Bcc").Where("message_id=?", imapEnvelope.MessageId).First(&envelope).Error

	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			envelope.Date      = imapEnvelope.Date
			envelope.InReplyTo = imapEnvelope.InReplyTo
			envelope.MessageId = imapEnvelope.MessageId
			envelope.Subject   = imapEnvelope.Subject

			for _, imapAddr := range imapEnvelope.From {
				addr, dbErr := GetOrCreateAddress(imapAddr)
				if dbErr != nil {
					err = dbErr
					return
				}
				envelope.From = append(envelope.From, addr)
			}

			for _, imapAddr := range imapEnvelope.Sender {
				addr, dbErr := GetOrCreateAddress(imapAddr)
				if dbErr != nil {
					err = dbErr
					return
				}
				envelope.Sender = append(envelope.Sender, addr)
			}

			for _, imapAddr := range imapEnvelope.ReplyTo {
				addr, dbErr := GetOrCreateAddress(imapAddr)
				if dbErr != nil {
					err = dbErr
					return
				}
				envelope.ReplyTo = append(envelope.ReplyTo, addr)
			}

			for _, imapAddr := range imapEnvelope.To {
				addr, dbErr := GetOrCreateAddress(imapAddr)
				if dbErr != nil {
					err = dbErr
					return
				}
				envelope.To = append(envelope.To, addr)
			}

			for _, imapAddr := range imapEnvelope.Cc {
				addr, dbErr := GetOrCreateAddress(imapAddr)
				if dbErr != nil {
					err = dbErr
					return
				}
				envelope.Cc = append(envelope.Cc, addr)
			}

			for _, imapAddr := range imapEnvelope.Bcc {
				addr, dbErr := GetOrCreateAddress(imapAddr)
				if dbErr != nil {
					err = dbErr
					return
				}
				envelope.Bcc = append(envelope.Bcc, addr)
			}

			db.Create(&envelope)

			PutEnvelopeInSearch(&envelope)
		} else {
			err = dbErr
		}
	}

	return
}

func PutEnvelopeInSearch(e *Envelope) (err error) {
	newDoc := esEnvelope{MessageId: e.MessageId, InReplyTo: e.InReplyTo, Date: e.Date, Subject: e.Subject}

	for _, from := range e.From {
		newDoc.From = append(newDoc.From, fmt.Sprintf("\"%s\" <%s@%s>", from.PersonalName, from.MailboxName, from.HostName))
	}
	for _, sender := range e.Sender {
		newDoc.Sender = append(newDoc.Sender, fmt.Sprintf("\"%s\" <%s@%s>", sender.PersonalName, sender.MailboxName, sender.HostName))
	}
	for _, replyTo := range e.ReplyTo {
		newDoc.ReplyTo = append(newDoc.ReplyTo, fmt.Sprintf("\"%s\" <%s@%s>", replyTo.PersonalName, replyTo.MailboxName, replyTo.HostName))
	}
	for _, to := range e.To {
		newDoc.To = append(newDoc.To, fmt.Sprintf("\"%s\" <%s@%s>", to.PersonalName, to.MailboxName, to.HostName))
	}
	for _, cc := range e.Cc {
		newDoc.Cc = append(newDoc.Cc, fmt.Sprintf("\"%s\" <%s@%s>", cc.PersonalName, cc.MailboxName, cc.HostName))
	}
	for _, bcc := range e.Bcc {
		newDoc.Bcc = append(newDoc.Bcc, fmt.Sprintf("\"%s\" <%s@%s>", bcc.PersonalName, bcc.MailboxName, bcc.HostName))
	}

	put, err := es.Index().
		Index("mail_envelope").
		Type("doc").
		Id(strconv.Itoa(int(e.ID))).
		BodyJson(newDoc).
		Do(context.Background())
	if err != nil {
		logger.Errorf("Coud not create index 'mail_envelope': %s", err)
		return
	}
	logger.Debugf("Indexed envelope %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	logger.Debugf("  Subject: %s\n", newDoc.Subject)
	logger.Debugf("  From: %s, Sender: %s, ReplyTo: %s\n", newDoc.From, newDoc.Sender, newDoc.ReplyTo)
	logger.Debugf("  To: %s, Cc: %s, Bcc: %s\n", newDoc.To, newDoc.Cc, newDoc.Bcc)

	return
}
