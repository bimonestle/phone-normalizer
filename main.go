package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "bimo"
	password = "your-password"
	dbname   = "phone_normalizer"
)

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	// Connect / Open to database
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	err = createDB(db, dbname)
	if err != nil {
		panic(err)
	}
	db.Close()
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

// To format the inputted phone number
// into numbers format only using regexp. No other characters e.g: (),-,_
func normalize(phone string) string {
	// re := regexp.MustCompile("[^0-9]") // We only want characters between 0 & 9
	re := regexp.MustCompile("[\\D]") // Match any non digits
	return re.ReplaceAllString(phone, "")
}

// To format the inputted phone number
// into numbers format only. No other characters e.g: (),-,_
// func normalize(phone string) string {
// 	// The Correct format - 0123456789
// 	// It contains numbers only

// 	var buf bytes.Buffer

// 	// When string is iterated individually, it will output rune. Not string
// 	for _, ch := range phone {
// 		// If the string contains between these runes
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch) // write rune into the Buffer
// 		}
// 	}
// 	return buf.String() // convert it back into string
// }
