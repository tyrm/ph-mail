package models

import (
	"time"

	"github.com/emersion/go-imap"
	"github.com/jinzhu/gorm"
)

// A message.
type Envelope struct {
	gorm.Model

	MessageId string    `gorm:"unique;index;not null"`
	Date      time.Time
	Subject   string

	From      []Address `gorm:"many2many:mail_efroms"`
	Sender    []Address `gorm:"many2many:mail_esenders"`
	ReplyTo   []Address `gorm:"many2many:mail_ereplytos"`
	To        []Address `gorm:"many2many:mail_etos"`
	Cc        []Address `gorm:"many2many:mail_eccs"`
	Bcc       []Address `gorm:"many2many:mail_ebccs"`

	InReplyTo string
}

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
		} else {
			err = dbErr
		}
	}

	return
}