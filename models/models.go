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
	MailboxName string
	// The host name.
	HostName string
}

// A message.
type Envelope struct {
	MessageId string     `gorm:"primary_key"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Date      time.Time
	Subject   string

	From      []*Address `gorm:"many2many:mail_efroms"`
	Sender    []*Address `gorm:"many2many:mail_esenders"`
	ReplyTo   []*Address `gorm:"many2many:mail_ereplytos"`
	To        []*Address `gorm:"many2many:mail_etos"`
	Cc        []*Address `gorm:"many2many:mail_eccs"`
	Bcc       []*Address `gorm:"many2many:mail_ebccs"`

	InReplyTo *Envelope
}