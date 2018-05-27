package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// An address.
type Address struct {
	gorm.Model

	// The personal name.
	PersonalName string
	// The SMTP at-domain-list (source route).
	AtDomainList string
	// The mailbox name.
	MailboxName string `gorm:"size:64"`
	// The host name.
	HostName string `gorm:"size:255"`
}

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