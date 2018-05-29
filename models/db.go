package models

import (
	"fmt"
	"log"
	"regexp"

	"../util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func CloseDB() {
	db.Close()

	return
}

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

func InitDB(connectionString string) {
	var err error
	dialect, dbArgs := DecodeEngine(connectionString)
	db, err = gorm.Open(dialect, dbArgs)
	util.PanicOnError(err, "Coud not connect to database")

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return "mail_" + defaultTableName;
	}

	db.AutoMigrate(&Address{}, &Envelope{})
	// Create Index to Speed searching for addresses
	db.Model(&Address{}).AddIndex("idx_host_name_mailbox_name", "lower(host_name)", "lower(mailbox_name)", "deleted_at")
	db.Model(&Envelope{}).AddIndex("idx_message_id", "message_id", "deleted_at")

	log.Printf("Connected to %s database", dialect)

	return
}
