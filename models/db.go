package models

import (
	"fmt"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/juju/loggo"
	"github.com/olivere/elastic"
	"golang.org/x/net/context"
)

var db *gorm.DB
var es *elastic.Client
var logger *loggo.Logger

func CloseDB() {
	db.Close()
	es.Stop()

	return
}

func DecodeEngine(engine string) (dialect string, args string) {
	pgRe, err := regexp.Compile(`postgresql://([\w]*):([\w\-.~:/?#\[\]!$&'()*+,;=]*)@([\w.]*)/([\w]*)`)
	if err != nil {
		logger.Criticalf("Regex compile error: %s", err)
		panic("PANIC!")
	}

	if pgRe.MatchString(engine) {
		dialect = "postgres"
		match := pgRe.FindStringSubmatch(engine)
		args = fmt.Sprintf("host=%s user=%s dbname=%s password=%s", match[3], match[1], match[4], match[2])
	} else {
		logger.Criticalf("Could not parse DB_ENGINE: %s", err)
		panic("PANIC!")
	}

	return
}

func InitDB(connectionString string, esHost string) {
	newLogger :=  loggo.GetLogger("mail.models")
	logger = &newLogger

	var err error
	dialect, dbArgs := DecodeEngine(connectionString)
	db, err = gorm.Open(dialect, dbArgs)
	if err != nil {
		logger.Criticalf("Coud not connect to database: %s", err)
		panic(err)
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return "mail_" + defaultTableName;
	}

	db.AutoMigrate(&Address{}, &Envelope{})
	// Create Index to Speed searching for addresses
	db.Model(&Address{}).AddIndex("idx_host_name_mailbox_name", "lower(host_name)", "lower(mailbox_name)", "deleted_at")
	db.Model(&Envelope{}).AddIndex("idx_message_id", "message_id", "deleted_at")

	logger.Infof("Connected to %s database", dialect)

	logger.Debugf("Trying to connecto to ElasticSearch: %s", esHost)
	es, err = elastic.NewClient(elastic.SetURL("http://10.1.68.60:9200"),  elastic.SetSniff(false))
	if err != nil {
		logger.Criticalf("Coud not connect ElasticSearch cluster: %s", err)
		panic(err)
	}

	// Check if the index called "mail_envelope" exists
	exists, err := es.IndexExists("mail_envelope").Do(context.Background())
	if err != nil {logger.Errorf("Error checking for: %s", err)}
	if !exists {
		logger.Infof("Creating index 'mail_envelope'")
		createIndex, err := es.CreateIndex("mail_envelope").BodyString(ESDocEnvelope).Do(context.Background())
		if err != nil {
			logger.Errorf("Coud not create index 'mail_envelope': %s", err)
			panic(err)
		}
		if !createIndex.Acknowledged {
			logger.Warningf("Creation of index 'mail_envelope' was unacknowledged")
		}
	}

	return
}
