package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/bimonestle/go-exercise-projects/08.Phone-Number-Normalizer/phone/phonedb"

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
	reset := phonedb.Reset("postgres", psqlInfo, dbname)
	if reset != nil {
		log.Fatal(reset)
	}

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	reset = phonedb.Migrate("postgres", psqlInfo)
	if reset != nil {
		log.Fatal(reset)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
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
	id, err := insertPhone(db, "(123) 456 7892")
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
	number, err := getPhone(db, id)
	if err != nil {
		log.Fatal("", err)
	}
	fmt.Println("ID is: ", id, ", Number is: ", number) // testing purpose
	// fmt.Println("id= ", id) // testing purpose

	phones, err := allPhones(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range phones {
		fmt.Printf("Working on...%+v\n", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("Updating or removing...", number)
			existing, err := findPhone(db, number)
			if err != nil {
				log.Fatal(err)
			}

			if existing != nil {
				// delete this number
				deletePhone(db, p.id)
			} else {
				// update this number
				p.number = number
				updatePhone(db, p)
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	statement := `SELECT * FROM phone_numbers WHERE id=$1`
	row := db.QueryRow(statement, id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", nil
	}

	return number, nil
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	fmt.Println("findPhone()")
	var p phone
	statement := `SELECT * FROM phone_numbers WHERE value=$1`
	row := db.QueryRow(statement, number)
	err := row.Scan(&p.id, &p.number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func updatePhone(db *sql.DB, p phone) error {
	fmt.Println("updatePhone()")
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.Exec(statement, p.id, p.number)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	fmt.Println("deletePhone()")
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.Exec(statement, id)
	return err
}

type phone struct {
	id     int
	number string
}

func allPhones(db *sql.DB) ([]phone, error) {
	statement := `SELECT id, value FROM phone_numbers`
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
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
