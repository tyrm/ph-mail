package models

import (
	"github.com/emersion/go-imap"
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

func GetOrCreateAddress(imapAddress *imap.Address) (address Address, err error) {
	dbErr := db.Where("LOWER(mailbox_name)=LOWER(?) AND LOWER(host_name)=LOWER(?)", imapAddress.MailboxName, imapAddress.HostName).First(&address).Error
	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			address.PersonalName = imapAddress.PersonalName
			address.AtDomainList = imapAddress.AtDomainList
			address.MailboxName  = imapAddress.MailboxName
			address.HostName     = imapAddress.HostName

			db.Create(&address)
		} else {
			err = dbErr
		}
	}

	return
}
