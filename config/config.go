package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"../util"
)

type Config struct {
	DBEngine       string

	MailUsername   string
	MailPassword   string
	MailIMAPServer string
}

func CollectConfig() (config Config) {
	var missingEnv []string

	// DB_ENGINE
	config.DBEngine = os.Getenv("DB_ENGINE")
	if config.DBEngine == "" {
		missingEnv = append(missingEnv, "DB_ENGINE")
	}

	// MAIL_USERNAME
	config.MailUsername = os.Getenv("MAIL_USERNAME")
	if config.MailUsername == "" {
		missingEnv = append(missingEnv, "MAIL_USERNAME")
	}

	// MAIL_PASSWORD
	config.MailPassword = os.Getenv("MAIL_PASSWORD")
	if config.MailPassword == "" {
		missingEnv = append(missingEnv, "MAIL_PASSWORD")
	}

	// IMAP_SERVER
	config.MailIMAPServer = os.Getenv("IMAP_SERVER")
	if config.MailIMAPServer == "" {
		missingEnv = append(missingEnv, "IMAP_SERVER")
	}

	// Validation
	if len(missingEnv) > 0 {
		var msg string = fmt.Sprintf("Environment variables missing: %v", missingEnv)
		log.Fatal(msg)
		panic(fmt.Sprint(msg))
	}

	return
}

func DecodeEngine(engine string) (dialect string, args string) {
	pgRe, err := regexp.Compile(`postgresql://([\w]*):([\w\-.~:/?#\[\]!$&'()*+,;=]*)@([\w.]*)/([\w]*)`)
	util.PanicOnError(err, "Regex compile error")

	if pgRe.MatchString(engine) {
		dialect = "postgres"
		match := pgRe.FindStringSubmatch(engine)
		args = fmt.Sprintf("host=%s user=%s dbname=%s password=%s", match[3], match[1], match[4], match[2])
	} else {
		panic(fmt.Sprint("Could not parse DB_ENGINE"))
	}

	return
}