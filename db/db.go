package db

import (
	"fmt"
	"log"
	"regexp"

	"../models"
	"../util"
	"github.com/emersion/go-imap"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DecodeEngine(engine string) (dialect string, args string) {
	pgRe, err := regexp.Compile(`postgresql://([\w]*):([\w\-.~:/?#\[\]!$&'()*+,;=]*)@([\w.]*)/([\w]*)`)
	if err != nil {
		log.Fatalf("Regex compile error: %s", err)
		panic(fmt.Sprintf("Regex compile error: %s", err))
	}

	if pgRe.MatchString(engine) {
		dialect = "postgres"
		match := pgRe.FindStringSubmatch(engine)
		args = fmt.Sprintf("host=%s user=%s dbname=%s password=%s", match[3], match[1], match[4], match[2])
	} else {
		panic(fmt.Sprint("Could not parse DB_ENGINE"))
	}

	return
}

func GetClient(connectionString string) *gorm.DB {
	dialect, dbArgs := DecodeEngine(connectionString)
	db, err := gorm.Open(dialect, dbArgs)
	util.PanicOnError(err, "Coud not connect to database")

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return "mail_" + defaultTableName;
	}
	db.AutoMigrate(&models.Address{}, &models.Envelope{})
	log.Printf("Connected to %s database", dialect)

	return db
}

func GetOrCreateAddress(dbClient *gorm.DB, imapAddress *imap.Address) (address models.Address, err error) {
	dbErr := dbClient.Where("LOWER(mailbox_name)=LOWER(?) AND LOWER(host_name)=LOWER(?)", imapAddress.MailboxName, imapAddress.HostName).First(&address).Error
	log.Printf("Error: %s", dbErr)
	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			address.PersonalName = imapAddress.PersonalName
			address.AtDomainList = imapAddress.AtDomainList
			address.MailboxName  = imapAddress.MailboxName
			address.HostName     = imapAddress.HostName

			log.Println(address.AtDomainList)

			dbClient.Create(&address)
		} else {
			err = dbErr
		}
	}

	return
}