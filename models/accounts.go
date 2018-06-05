package models

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model

	Address string
	Envelope []*Envelope
}