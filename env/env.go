package env

import (
	"github.com/jinzhu/gorm"
	"github.com/emersion/go-imap/client"
)

type Env struct {
	Config Config
	DB     *gorm.DB
	IMAP   *client.Client
}