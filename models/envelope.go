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

func EnvelopeExistsByMsgID(id string) bool {
	var count int64
	db.Model(&Envelope{}).Where("message_id = ?", id).Count(&count)

	return count > 0
}

func CreateEnvelope(imapEnvelope *imap.Envelope) (envelope Envelope, err error) {
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

	IndexEnvelope(&envelope)

	return
}

func GetEnvelopeByMsgID(mid string) (envelope Envelope, err error) {
	err = db.Preload("From").Preload("Sender").Preload("ReplyTo").Preload("To").
		Preload("Cc").Preload("Bcc").Where("message_id=?", mid).First(&envelope).Error

	return
}

func IndexEnvelope(e *Envelope) (err error) {
	newDoc := e.ToES()

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

func (e Envelope) ToES() (es *esEnvelope) {
	es = &esEnvelope{
		MessageId: e.MessageId,
		InReplyTo: e.InReplyTo,
		Date: e.Date,
		Subject: e.Subject,
	}

	for _, from := range e.From {
		es.From = append(es.From, fmt.Sprintf("\"%s\" <%s@%s>", from.PersonalName, from.MailboxName, from.HostName))
	}
	for _, sender := range e.Sender {
		es.Sender = append(es.Sender, fmt.Sprintf("\"%s\" <%s@%s>", sender.PersonalName, sender.MailboxName, sender.HostName))
	}
	for _, replyTo := range e.ReplyTo {
		es.ReplyTo = append(es.ReplyTo, fmt.Sprintf("\"%s\" <%s@%s>", replyTo.PersonalName, replyTo.MailboxName, replyTo.HostName))
	}
	for _, to := range e.To {
		es.To = append(es.To, fmt.Sprintf("\"%s\" <%s@%s>", to.PersonalName, to.MailboxName, to.HostName))
	}
	for _, cc := range e.Cc {
		es.Cc = append(es.Cc, fmt.Sprintf("\"%s\" <%s@%s>", cc.PersonalName, cc.MailboxName, cc.HostName))
	}
	for _, bcc := range e.Bcc {
		es.Bcc = append(es.Bcc, fmt.Sprintf("\"%s\" <%s@%s>", bcc.PersonalName, bcc.MailboxName, bcc.HostName))
	}
	return
}