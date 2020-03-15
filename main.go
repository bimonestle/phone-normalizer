package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "rizalbimanto"
	password = "your-password"
	dbname   = "phone_normalizer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	// Connect / Open to database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		// panic(err)
		log.Fatal("Fatel log: ", err)
	}
	// err = createDB(db, dbname)
	// if err != nil {
	// 	panic(err)
	// }

	// If this part is commented because it has resetDB(db, dbname)
	// so that the id won't be incremented.
	err = resetDB(db, dbname)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createPhoneNumbTable(db)
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "1234567890")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "123 456 7891")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "(123) 456 7892")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "(123) 456-7893")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "123-456-7894")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "123-456-7890")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "1234567892")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(db, "(123)456-7892")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal("Error: ", err)
	}
	// fmt.Println("id= ", id) // testing purpose
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id;`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		// will return an invalid ID and the error
		return -1, err
	}
	return id, nil
}

func createPhoneNumbTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`
	_, err := db.Exec(statement)
	return err
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
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
