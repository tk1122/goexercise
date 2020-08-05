package main

import (
	"bytes"
	"fmt"
	dbpkg "goercises/phone/db"
	"log"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "gophercise_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(dbpkg.Reset("postgres", psqlInfo, dbname))
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(dbpkg.Migrate("postgres", psqlInfo))
	db, err := dbpkg.Open("postgres", psqlInfo)
	must(err)
	must(db.Seed())

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		_ = p
		normalizedNumber := normalizeRegex(p.Number)
		existing, err := db.FindPhoneExcept(normalizedNumber, p.Id)
		must(err)
		if existing != nil {
			must(db.DeletePhone(p.Id))
		} else {
			p.Number = normalizedNumber
			must(db.UpdatePhone(p))
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func normalizeRegex(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}
